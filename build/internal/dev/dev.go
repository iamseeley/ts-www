package dev

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"ts-www/build/internal/config"
	"ts-www/build/internal/models"
	"ts-www/build/internal/ogimage"
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
					createTemplateForDir(filepath.Base(event.Name), templateDir)
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
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
						appendFrontmatter(event.Name, collection)
					}
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func appendFrontmatter(filePath, collection string) error {
	// Define frontmatter templates for each collection
	templates := map[string]string{
		"notes": `---
title: Your Note Title
summary: A brief summary of the note.
tags: [tag1, tag2]
date: YYYY-MM-DD
draft: false
---
`,
		"logs": `---
title: Your Log Title
date: YYYY-MM-DD
draft: false
content: 
---
`,
		"page": `---
title: Your Page Title
description: A brief description of the page
draft: false
---
`,
		"posts": `---
title: Your Post Title
description: A brief description of the post
date: YYYY-MM-DD
draft: false
---
`,
		"collections": `---
title: Your Collection Title
description: A brief description of the collection
type: Link, Book, Blog, etc...?
draft: false
---
`,
		// Add more templates for other collections as needed
	}

	// Select the appropriate template based on the collection name
	template, ok := templates[collection]
	if !ok {
		log.Printf("No frontmatter template for collection: %s", collection)
		return fmt.Errorf("no frontmatter template for collection: %s", collection)
	}

	// Read the existing content of the file
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

	// Write the template and then the existing content to the file
	if _, err := file.WriteString(template + string(existingContent)); err != nil {
		return err
	}

	log.Printf("Frontmatter appended to file in collection '%s': %s", collection, filePath)
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

	p, err := utils.LoadPageFromDirectory(cfg.ContentPath, filePath)
	if err != nil {
		log.Printf("Error loading page: %v", err) // Log the error for debugging
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	collection := filepath.Base(filepath.Dir(filePath))

	templateData := struct {
		Page *models.Content
		Data map[string]interface{}
	}{
		Page: p,
		Data: data,
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

	go ogimage.GenerateAllOGImages(cfg.ContentPath, "assets/og-image/")
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

	// http.Handle("/", http.RedirectHandler("/index", http.StatusSeeOther))
	fs := http.FileServer(http.Dir("src/public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	// http.HandleFunc("/", makeHandler(pageHandler))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
