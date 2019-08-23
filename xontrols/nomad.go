package xontrols

import "github.com/louisevanderlith/droxolite/context"

type NomadController interface {
	Get(context.Contexer) error
}
