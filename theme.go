package droxolite

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/husk"
)

//GetDefaultTheme will attempt to contact the Theme.API for the profile's site
func GetDefaultTheme(host, instanceID, siteName string) (bodies.ThemeSetting, error) {
	prof := bodies.ThemeSetting{}
	_, err := DoGET("", &prof, instanceID, "Folio.API", "theme", siteName)

	if err != nil {
		return bodies.ThemeSetting{}, err
	}

	return prof, nil
}

//GetNoTheme returns an empty Theme. This should be used for when you don't need a full theme.
func GetNoTheme(host, instanceID, siteName string) bodies.ThemeSetting {
	return bodies.NewThemeSetting(siteName, host, husk.CrazyKey(), instanceID, "UA-000000000-0")
}

//UpdateTheme downloads the latest master templates from Theme.API
func UpdateTheme(instanceID string) error {
	lst, err := findTemplates(instanceID)

	if err != nil {
		return err
	}

	url, err := GetServiceURL(instanceID, "Theme.API", false)

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
	code, err := DoGET("", &result, instanceID, "Theme.API", "asset", "html")

	if err != nil {
		return []string{strconv.Itoa(code)}, err
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
