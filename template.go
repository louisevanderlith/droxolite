package droxolite

import (
	"encoding/json"
	"fmt"
	"github.com/louisevanderlith/kong/tokens"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

//UpdateTemplate downloads the latest master templates from Theme.API
func UpdateTemplate(access string, claims tokens.Claimer) error {
	url, err := claims.GetResourceURL("theme")

	if err != nil {
		return err
	}

	lst, err := findTemplates(access, url)

	if err != nil {
		return err
	}

	for _, v := range lst {
		err = downloadTemplate(access, v, url)

		if err != nil {
			return err
		}
	}

	return nil
}

func findTemplates(access, themeUrl string) ([]string, error) {
	fullURL := fmt.Sprintf("%s/asset/html", themeUrl)

	req := httptest.NewRequest(http.MethodGet, fullURL, nil)
	req.Header.Set("Authorization", "Bearer "+access)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result []string
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func downloadTemplate(access, template, themeURL string) error {
	fullURL := fmt.Sprintf("%s/asset/html/%s", themeURL, template)
	req := httptest.NewRequest(http.MethodGet, fullURL, nil)
	req.Header.Set("Authorization", "Bearer "+access)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create("/views/_shared/" + template)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

//LoadTemplate parses html files, to be used as layouts
func LoadTemplate(viewPath, masterpage string) (*template.Template, *template.Template, error) {
	temps, err := template.ParseFiles(findFiles(viewPath)...)

	if err != nil {
		return nil, nil, err
	}

	return template.New(masterpage), temps, nil
}

func findFiles(templatesPath string) []string {
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
