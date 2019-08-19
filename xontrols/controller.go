package xontrols

import (
	"io"

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
	Ctx() context.Contexer
}

type APIController interface {
	Controller
	ServeBinary(data []byte, filename string) error
	ServeBinaryWithMIME(data []byte, filename, mimetype string) error
	ServeBinaryStream(data io.Reader, filename, mimetype string) error
}

type UIController interface {
	APIController
	SetTheme(settings bodies.ThemeSetting, masterpage string)
	CreateSideMenu(menu *bodies.Menu)
}
