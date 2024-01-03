package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

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

// func FilterNotesLastSixMonths(notes []models.Page) []models.Page {
//     filteredNotes := []models.Page{}
//     sixMonthsAgo := time.Now().AddDate(0, -6, 0) // Six months ago

//     for _, note := range notes {
//         if dateStr, ok := note.FrontMatter["date"].(string); ok {
//             noteDate, err := time.Parse("2006-01-02", dateStr)
//             if err == nil && noteDate.After(sixMonthsAgo) {
//                 filteredNotes = append(filteredNotes, note)
//             }
//         }
//     }
//     return filteredNotes
// }
