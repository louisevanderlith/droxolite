package xontrols

import "github.com/louisevanderlith/droxolite/context"

//NomadController is the simplest form of controller. [GET]
type Nomad interface {
	Get(context.Requester) (int, interface{})
}
