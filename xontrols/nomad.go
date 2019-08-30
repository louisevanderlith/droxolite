package xontrols

import "github.com/louisevanderlith/droxolite/context"

//NomadController is the simplest form of controller.
type NomadController interface {
	Get(context.Contexer) error
}
