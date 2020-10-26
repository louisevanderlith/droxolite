package drx

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//LoadTemplate parses html files, to be used as layouts
func LoadTemplate(viewPath string) (*template.Template, error) {
	return template.ParseFiles(FindFiles(viewPath)...)
}

func FindFiles(templatesPath string) []string {
	var result []string

	filepath.Walk(templatesPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		}

		if !f.IsDir() && strings.HasSuffix(path, ".html") {
			result = append(result, path)
		}

		return nil
	})

	return result
}
