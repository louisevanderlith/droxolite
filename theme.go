package droxolite

import (
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

	prof.InstanceID = instanceID
	prof.Host = host
	return prof, nil
}

//GetNoTheme returns an empty Theme. This should be used for when you don't need a full theme.
func GetNoTheme(host, instanceID, siteName string) bodies.ThemeSetting {
	return bodies.NewThemeSetting(siteName, host, husk.CrazyKey(), instanceID, "UA-000000000-0")
}
