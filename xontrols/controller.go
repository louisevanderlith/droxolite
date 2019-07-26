package xontrols

import (
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
)

//Controller provides the interface for all controllers in droxolite
type Controller interface {
	Filter() bool
	CreateInstance(ctx context.Contexer)
	Prepare()
	Serve(int, error, interface{}) error
}

type APIController interface {
	Controller
}

type UIController interface {
	APIController
	SetTheme(settings bodies.ThemeSetting)
}
