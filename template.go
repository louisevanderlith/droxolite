package droxolite

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/louisevanderlith/droxolite/do"
)

//UpdateTheme downloads the latest master templates from Theme.API
func UpdateTheme(instanceID string) error {
	lst, err := findTemplates(instanceID)

	if err != nil {
		return err
	}

	url, err := do.GetServiceURL(instanceID, "Theme.API", false)

	if err != nil {
		return err
	}

	for _, v := range lst {
		err = downloadTemplate(instanceID, v, url)

		if err != nil {
			return err
		}
	}

	return nil
}

func findTemplates(instanceID string) ([]string, error) {
	result := []string{}
	_, err := do.GET("", &result, instanceID, "Theme.API", "asset", "html")

	if err != nil {
		return []string{}, err
	}

	return result, nil
}

func downloadTemplate(instanceID, template, themeURL string) error {
	fullURL := fmt.Sprintf("%sasset/html/%s", themeURL, template)
	resp, err := http.Get(fullURL)

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
