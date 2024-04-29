package dev

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"ts-www/build/internal/config"
	"ts-www/build/internal/models"
	"ts-www/build/internal/utils"

	"github.com/fsnotify/fsnotify"
)

func watchContentDirectory(contentDir, templateDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(contentDir)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					// Create default template for the new directory
					createTemplateForDir(filepath.Base(event.Name), templateDir)

					// Create base.md in the new directory
					createBaseMD(event.Name)
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func createBaseMD(dirPath string) {
	baseMDContent := `---
title: "Title"
description: "Description"
date: "YYYY-MM-DD"
draft: true
---
`
	baseMDPath := filepath.Join(dirPath, "base.md")
	err := os.WriteFile(baseMDPath, []byte(baseMDContent), 0644)
	if err != nil {
		log.Printf("Failed to create base.md in %s: %v", dirPath, err)
	} else {
		log.Printf("Created base.md in %s", dirPath)
	}
}

func watchForNewMarkdownFiles(contentDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add contentDir and its immediate subdirectories to the watcher
	err = filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && (path == contentDir || filepath.Dir(path) == contentDir) {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err != nil || info.IsDir() {
					continue
				}
				if filepath.Ext(event.Name) == ".md" {
					collection := filepath.Base(filepath.Dir(event.Name))
					if collection != filepath.Base(contentDir) { // Exclude files directly in contentDir
						appendFrontmatter(event.Name, collection, contentDir)
					}
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func appendFrontmatter(filePath, collection string, contentDir string) error {
	// Construct the path for base.md in the collection
	baseFilePath := filepath.Join(contentDir, collection, "base.md")

	// Read the content of base.md
	baseContent, err := os.ReadFile(baseFilePath)
	if err != nil {
		log.Printf("Error reading base.md for collection %s: %v", collection, err)
		return err
	}

	// Extract frontmatter from base.md
	start := bytes.Index(baseContent, []byte("---"))
	if start == -1 {
		return fmt.Errorf("frontmatter delimiter not found in base.md of collection: %s", collection)
	}
	end := bytes.Index(baseContent[start+3:], []byte("---"))
	if end == -1 {
		return fmt.Errorf("closing frontmatter delimiter not found in base.md of collection: %s", collection)
	}
	frontmatterTemplate := string(baseContent[start : start+3+end+3])

	// Read the existing content of the file being processed
	existingContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Create a new file with the same name (overwriting the existing file)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the frontmatter template and then the existing content to the file
	if _, err := file.WriteString(frontmatterTemplate + "\n" + string(existingContent)); err != nil {
		return err
	}

	log.Printf("Frontmatter from base.md appended to file in collection '%s': %s", collection, filePath)
	return nil
}

func createTemplateForDir(dirName, templateDir string) {
	templatePath := filepath.Join(templateDir, dirName+".html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// Basic HTML template content
		templateContent := fmt.Sprintf(`
{{ define "%s" }}
{{template "_top" .}}

	<section>
		<h2>{{.Page.Title}}</h2>
		<article>
		{{ .Page.Body | markDown }}
		</article>
	</section>

{{template "_bottom" .}}
{{ end }}
`, dirName)

		err = os.WriteFile(templatePath, []byte(templateContent), 0644)
		if err != nil {
			log.Printf("Failed to create template for %s: %v", dirName, err)
		} else {
			log.Printf("Created template: %s", templatePath)
		}
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request, filePath string) {
	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	setCacheHeaders(w, 600)

	log.Printf("Constructed file path: %s", filePath)

	data, err := utils.LoadData(cfg.DataPath)
	if err != nil {
		log.Printf("Failed to load data: %v", err)
	}

	feed, err := utils.LoadFeed(cfg.ContentPath)
	if err != nil {
		log.Println("Failed to load feed:", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	p, err := utils.LoadPageFromDirectory(cfg.ContentPath, filePath)
	if err != nil {
		log.Printf("Error loading page: %v", err) // Log the error for debugging
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	// Generate the OG Image URL
	// ogImageFileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)) + "-og-image.png"
	// ogImageUrl := "/public/og-image/" + ogImageFileName
	// p.OGImageURL = ogImageUrl

	collection := filepath.Base(filepath.Dir(filePath))

	templateData := struct {
		Page *models.Content
		Data map[string]interface{}
		Feed []models.Content
	}{
		Page: p,
		Data: data,
		Feed: feed,
	}

	tmplName := collection
	tmpl := utils.Templates.Lookup(tmplName + ".html")
	if tmpl == nil {
		log.Printf("Template %s.html not found, using default site.html", tmplName)
		tmpl = utils.Templates.Lookup("site.html")
	}

	utils.RenderTemplateDev(w, tmplName, templateData)
}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[1])
	}
}

func setCacheHeaders(w http.ResponseWriter, maxAge int) {
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
}

func StartServer() {
	// Load configuration
	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = utils.ConvertMarkdownToJSON(cfg.ContentPath, cfg.DataPath)
	if err != nil {
		log.Fatalf("Error converting markdown to JSON: %v", err)
	}

	// Copy the theme CSS to the assets/css directory
	themeCSSPath := filepath.Join("themes", cfg.ThemeName+".css")
	assetsCSSPath := filepath.Join("assets/css", cfg.ThemeName+".css")
	os.MkdirAll(filepath.Dir(assetsCSSPath), os.ModePerm) // Create the assets/css directory
	err = utils.CopyFile(themeCSSPath, assetsCSSPath)
	if err != nil {
		log.Fatalf("Failed to copy theme CSS to assets directory: %v", err)
	}

	err = utils.LoadTemplates()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	outputDir := cfg.OutputPath

	os.MkdirAll(outputDir, os.ModePerm)

	// Copy the assets directory to public in the output directory
	assetsSrc := "assets"
	assetsDst := filepath.Join(outputDir, "public")
	err = utils.CopyDir(assetsSrc, assetsDst)
	if err != nil {
		log.Fatalf("Failed to copy assets directory: %v", err)
	}

	// go ogimage.GenerateAllOGImages(cfg.ContentPath, "assets/og-image/")
	go watchContentDirectory("content", "templates")
	go watchForNewMarkdownFiles("content")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/":
			// Serve 'index.md' from the 'page' directory for the root path
			pageHandler(w, r, "page/index.md")
		case !strings.Contains(r.URL.Path[1:], "/"):
			// Handle 'page' collection routes without the 'page' prefix in the URL
			// For example, "/about" will serve "page/about.md"
			strippedPath := strings.TrimPrefix(r.URL.Path, "/")
			pageHandler(w, r, fmt.Sprintf("page/%s.md", strippedPath))
		default:
			// Handle other collection routes
			// For example, "/post/post1" will serve "post/post1.md"
			pageHandler(w, r, r.URL.Path[1:]+".md") // [1:] to remove the leading '/'
		}
	})

	fs := http.FileServer(http.Dir("src/public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
