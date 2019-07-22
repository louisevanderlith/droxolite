package droxolite

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/louisevanderlith/droxolite/bodies"
	folio "github.com/louisevanderlith/folio/core"
)

func GetDefaultTheme(host, instanceID, siteName string) (bodies.ThemeSetting, error) {
	prof := folio.Profile{}
	_, err := DoGET("", &prof, instanceID, "Folio.API", "profile", siteName)

	if err != nil {
		return bodies.ThemeSetting{}, err
	}

	result := bodies.NewThemeSetting(prof.Title, host, prof.ImageKey, instanceID, prof.GTag)

	return result, nil
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
	fullURL := fmt.Sprintf("%sv1/%s/%s/%s", themeURL, "asset", "html", template)
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
