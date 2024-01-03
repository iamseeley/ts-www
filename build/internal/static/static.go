package static

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"ts-www/build/internal/config"
	"ts-www/build/internal/models"
	"ts-www/build/internal/utils"

	"github.com/russross/blackfriday/v2"
)

// Load site pages written in Markdown from a directory
func loadPageFromDirectory(directory, title string) (*models.Page, error) {
	filename := directory + title
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	frontMatter, body, err := utils.ParseFrontMatter(content)
	if err != nil {
		return nil, err
	}

	var page models.Page
	if title, ok := frontMatter["title"].(string); ok {
		page.Title = title
	}
	if description, ok := frontMatter["description"].(string); ok {
		page.Description = description
	}

	cfg, err := config.LoadConfig("./config.json") // Load configuration
	if err != nil {
		return nil, err
	}

	page.Theme = cfg.ThemeName
	page.Body = blackfriday.Run(body)
	page.Collection = filepath.Base(filepath.Dir(filename))

	return &page, nil
}

func loadData(directory string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			fileData, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var jsonData interface{}
			if err := json.Unmarshal(fileData, &jsonData); err != nil {
				return err
			}

			key := filepath.Base(path)
			key = strings.TrimSuffix(key, filepath.Ext(key)) // Use filename as the key
			data[key] = jsonData
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func loadTemplates() error {
	var err error
	templates, err = template.New("").Funcs(template.FuncMap{"markDown": markDowner}).ParseGlob("templates/*.html")
	if err != nil {
		return fmt.Errorf("error loading templates: %w", err)
	}
	return nil
}

func init() {
	err := loadTemplates()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
}

// BuildSite generates static HTML files from Markdown content
func BuildSite() {
	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	data, err := loadData("data")
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Copy the theme CSS to the assets/css directory
	themeCSSPath := filepath.Join("themes", cfg.ThemeName+".css")
	assetsCSSPath := filepath.Join("assets/css", cfg.ThemeName+".css")
	os.MkdirAll(filepath.Dir(assetsCSSPath), os.ModePerm) // Create the assets/css directory if it doesn't exist
	err = utils.CopyFile(themeCSSPath, assetsCSSPath)
	if err != nil {
		log.Fatalf("Failed to copy theme CSS to assets directory: %v", err)
	}

	contentDir := cfg.ContentPath
	outputDir := cfg.OutputPath

	os.MkdirAll(outputDir, os.ModePerm)

	// Copy the assets directory to public in the output directory
	assetsSrc := "assets"
	assetsDst := filepath.Join(outputDir, "public")
	err = utils.CopyDir(assetsSrc, assetsDst)
	if err != nil {
		log.Fatalf("Failed to copy assets directory: %v", err)
	}

	// ... rest of the BuildSite function ...

	// Iterate through all Markdown files in the content directory
	err = filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			return generateHTML(path, outputDir, data, cfg) // Pass the loaded data and cfg
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error building site: %v", err)
	}

	log.Println("Site built successfully")
}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

func renderTemplate(tmpl string, content interface{}) {
	err := templates.ExecuteTemplate(os.Stdout, tmpl, content)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

var templates *template.Template

func generateHTML(mdPath, outputDir string, data map[string]interface{}, cfg *config.Config) error {
	// Extract the relative path of the Markdown file from the content directory
	relativePath, err := filepath.Rel(cfg.ContentPath, mdPath)
	if err != nil {
		log.Printf("Error getting relative path: %v", err)
		return err
	}

	// Change the extension from .md to .html
	htmlPath := strings.TrimSuffix(relativePath, filepath.Ext(relativePath)) + ".html"

	// If the file is in the 'page' directory, place it at the root of the output directory
	if strings.HasPrefix(relativePath, "page"+string(filepath.Separator)) {
		htmlPath = strings.TrimPrefix(htmlPath, "page"+string(filepath.Separator))
	}

	// Create the full path for the output HTML file
	outputPath := filepath.Join(outputDir, htmlPath)

	// Create the necessary directories in the output path
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		log.Printf("Error creating directories: %v", err)
		return err
	}

	// page, err := loadPageFromDirectory(filepath.Dir(mdPath)+"/", filepath.Base(mdPath))
	// if err != nil {
	// 	log.Printf("Error loading page from directory: %v", err)
	// 	return err
	// }

	page, err := loadPageFromDirectory(filepath.Dir(mdPath)+"/", filepath.Base(mdPath))
	if err != nil {
		log.Printf("Error loading page: %v", err)

		return err // Log the error for debugging
	}

	// Determine template based on the collection (parent directory name)
	collection := filepath.Base(filepath.Dir(mdPath))
	tmplName := collection

	// Use the collection's template; default to "site.html" if not found
	tmpl := templates.Lookup(tmplName + ".html")
	if tmpl == nil {
		log.Printf("Template %s.html not found, using default site.html", tmplName)
		tmpl = templates.Lookup("site.html")
	}

	newdata, err := loadData("data")
	if err != nil {
		log.Printf("Failed to load data: %v", err)
	}

	log.Printf("Executing template with Page: %+v", page)

	templateData := struct {
		Page *models.Page
		Data map[string]interface{}
	}{
		Page: page,
		Data: newdata,
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer outputFile.Close()

	err = tmpl.ExecuteTemplate(outputFile, tmplName, templateData)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		return err
	}
	// Save the rendered HTML
	// err = os.WriteFile(outputPath, templateData.Page.Body, 0644)
	// if err != nil {
	// 	log.Printf("Error writing file: %v", err)
	// 	return err
	// }

	return nil
}
