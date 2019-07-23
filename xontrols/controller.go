package xontrols

import (
	"github.com/louisevanderlith/droxolite/context"
)

//Controller provides the interface for all controllers in droxolite
type Controller interface {
	Filter() bool
	CreateInstance(ctx context.Contexer)
	Prepare()
}

type APIController interface {
	Controller
}

type UIController interface {
	APIController
	Render()
}
