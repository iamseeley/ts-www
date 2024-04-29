package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"ts-www/build/internal/config"
	"ts-www/build/internal/models"

	"github.com/go-yaml/yaml"
	"github.com/russross/blackfriday/v2"
)

func LoadFeed(directory string) ([]models.Content, error) {
	var allContent []models.Content

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil // Skip directories and non-markdown files.
		}

		// Extract collection name and file name.
		collectionName := filepath.Base(filepath.Dir(path))
		fileName := filepath.Base(path)

		// Skip if the collection is "page" or the file is "base.md".
		if collectionName == "page" || fileName == "base.md" {
			return nil
		}

		// Correctly construct the path to include the collection
		collectionDir := filepath.Dir(path)

		// Use LoadPageFromDirectory to load the content, passing the correct directory and file name.
		content, err := LoadPageFromDirectory(collectionDir+"/", fileName)
		if err != nil {
			log.Printf("Error loading content from %s: %v", path, err)
			return nil // Continue processing other files even if one fails.
		}

		// Append loaded content to the slice.
		allContent = append(allContent, *content)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort the slice by date.
	sort.Slice(allContent, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", allContent[i].Date)
		dateJ, errJ := time.Parse("2006-01-02", allContent[j].Date)
		if errI != nil || errJ != nil {
			log.Printf("Error parsing date: %v, %v", errI, errJ)
			return false
		}
		return dateI.After(dateJ)
	})

	return allContent, nil
}

var ErrDraftContent = errors.New("content is marked as draft")

func LoadPageFromDirectory(directory, title string) (*models.Content, error) {
	filename := directory + title
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	frontMatter, body, err := ParseFrontMatter(content)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadConfig("./config.json") // Load configuration
	if err != nil {
		return nil, err
	}

	var contentItem models.Content
	contentItem.Title, _ = frontMatter["title"].(string)
	contentItem.Date, _ = frontMatter["date"].(string)
	if description, ok := frontMatter["description"].(string); ok {
		contentItem.Description = description
	} else {
		contentItem.Description = ""
	}
	// Check if the content is marked as a draft
	if draft, ok := frontMatter["draft"].(bool); ok && draft {
		return nil, ErrDraftContent
	}
	if featured, ok := frontMatter["featured"].(bool); ok {
		contentItem.Featured = featured
	} else {
		contentItem.Draft = true
	}
	contentItem.Body = body
	contentItem.URL, _ = frontMatter["url"].(string)
	contentItem.Theme = cfg.ThemeName // Assuming the theme is consistent across all content
	contentItem.Collection = filepath.Base(filepath.Dir(filename))

	return &contentItem, nil
}

func LoadData(directory string) (map[string]interface{}, error) {
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

func MarkDowner(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

var Templates *template.Template

func LoadTemplates() error {
	var err error
	Templates, err = template.New("").Funcs(template.FuncMap{"markDown": MarkDowner, "parseDate": ParseDate, "now": Now}).ParseGlob("templates/*.html")
	if err != nil {
		return fmt.Errorf("error loading templates: %w", err)
	}
	return nil
}

func Init() {
	err := LoadTemplates()
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}
}

func RenderTemplateStatic(outputPath, tmpl string, content interface{}) error {
	// Create or open the file where the HTML will be saved
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outputFile.Close()

	// Execute the template and write the output to the file
	err = Templates.ExecuteTemplate(outputFile, tmpl, content)
	if err != nil {
		return fmt.Errorf("error rendering template: %v", err)
	}

	return nil
}

func RenderTemplateDev(w http.ResponseWriter, tmpl string, content interface{}) {
	err := Templates.ExecuteTemplate(w, tmpl, content)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseFrontMatter(content []byte) (map[string]interface{}, []byte, error) {
	frontMatter := make(map[string]interface{})
	var contentStart int

	delimiter := []byte("---")
	start := bytes.Index(content, delimiter)
	if start == -1 {
		return nil, nil, errors.New("Front matter delimiter not found")
	}

	end := bytes.Index(content[start+len(delimiter):], delimiter)
	if end == -1 {
		return nil, nil, errors.New("Second front matter delimiter not found")
	}

	if err := yaml.Unmarshal(content[start+len(delimiter):start+len(delimiter)+end], &frontMatter); err != nil {
		return nil, nil, err
	}

	contentStart = start + len(delimiter) + end + len(delimiter)
	actualContent := content[contentStart:]

	return frontMatter, actualContent, nil
}

func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func CopyDir(src string, dst string) error {
	// Get properties of source dir
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(src)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		srcFile := filepath.Join(src, obj.Name())
		dstFile := filepath.Join(dst, obj.Name())

		if obj.IsDir() {
			// Create sub-directories - recursively
			err = CopyDir(srcFile, dstFile)
			if err != nil {
				return err
			}
		} else {
			// Perform the file copy
			err = CopyFile(srcFile, dstFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ConvertMarkdownToJSON(contentDir, dataDir string) error {
	// Step 1: Create a set of current markdown filenames
	markdownFiles := make(map[string]struct{})
	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			relativePath, err := filepath.Rel(contentDir, path)
			if err != nil {
				return err
			}
			markdownFiles[relativePath] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Step 2: Process each collection and update JSON data
	collections, err := os.ReadDir(contentDir)
	if err != nil {
		return err
	}

	for _, collection := range collections {
		if collection.IsDir() {
			collectionName := collection.Name()
			jsonFileName := collectionName + ".json"
			jsonFilePath := filepath.Join(dataDir, jsonFileName)

			// Read existing JSON data
			existingData, err := os.ReadFile(jsonFilePath)
			if err != nil && !os.IsNotExist(err) {
				return err
			}

			dataMap := make(map[string]map[string]interface{})
			if len(existingData) > 0 {
				err = json.Unmarshal(existingData, &dataMap)
				if err != nil {
					return err
				}
			}

			// Update dataMap with new/updated markdown files
			for relativePath := range markdownFiles {
				if strings.HasPrefix(relativePath, collectionName+"/") {
					fullPath := filepath.Join(contentDir, relativePath)
					content, err := os.ReadFile(fullPath)
					if err != nil {
						return err
					}

					frontMatter, body, err := ParseFrontMatter(content)
					if err != nil {
						return err
					}

					fileIdentifier := strings.TrimSuffix(filepath.Base(fullPath), filepath.Ext(fullPath))
					jsonData := map[string]interface{}{
						"frontMatter": frontMatter,
						"body":        string(body),
					}

					dataMap[fileIdentifier] = jsonData
				}
			}

			// Check for deletions and remove entries from dataMap
			for key := range dataMap {
				relativePath := collectionName + "/" + key + ".md"
				if _, exists := markdownFiles[relativePath]; !exists {
					delete(dataMap, key)
				}
			}

			// Serialize and save the updated JSON data
			newJSONData, err := json.MarshalIndent(dataMap, "", "  ")
			if err != nil {
				return err
			}
			err = os.WriteFile(jsonFilePath, newJSONData, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ParseDate parses a date string in "YYYY-MM-DD" format and returns a time.Time object.
func ParseDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Handle the error according to your needs
		return time.Time{} // return zero time on error
	}
	return t
}

// Now returns the current time as a time.Time object.
func Now() time.Time {
	return time.Now()
}
