package bodies

import (
	"github.com/louisevanderlith/husk"
)

//ThemeSetting is the basic controls variables accessed by the Front-end
type ThemeSetting struct {
	LogoKey    husk.Key
	Name       string
	Host       string
	InstanceID string
	GTag       string
}

func NewThemeSetting(name, host string, logoKey husk.Key, instanceID, gtag string) ThemeSetting {
	return ThemeSetting{
		Name:       name,
		LogoKey:    logoKey,
		Host:       host, //os.Getenv("HOST")
		InstanceID: instanceID,
		GTag:       gtag,
	}
}
