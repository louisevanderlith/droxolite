package xontrols

import (
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/roletype"
)

//Controller provides the interface for all controllers in droxolite
type Controller interface {
	Filter(requiredRole roletype.Enum, publicKeyPath, serviceName string) bool
	CreateInstance(ctx context.Contexer, instanceID string)
	GetInstanceID() string
	Prepare()
	Serve(int, error, interface{}) error
}

type APIController interface {
	Controller
}

type UIController interface {
	APIController
	SetTheme(settings bodies.ThemeSetting, masterpage string)
	CreateSideMenu(menu *bodies.Menu)
}
