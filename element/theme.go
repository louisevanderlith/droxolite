package element

import (
	"github.com/louisevanderlith/droxolite/do"
	"github.com/louisevanderlith/husk"
)

//GetNoTheme returns an empty Theme. This should be used for when you don't need a full theme, ex. an API
func GetNoTheme(host, instanceID, siteName string) *Identity {
	return NewIdentity(siteName, host, husk.CrazyKey(), instanceID, "UA-000000000-0")
}

//GetDefaultTheme will attempt to contact the Theme.API for the profile's site
func GetDefaultTheme(host, instanceID, siteName string) (*Identity, error) {
	dentity := &Identity{}
	_, err := do.GET("", dentity, instanceID, "Folio.API", "theme", siteName)

	if err != nil {
		return nil, err
	}

	dentity.InstanceID = instanceID
	dentity.Host = host
	return dentity, nil
}
