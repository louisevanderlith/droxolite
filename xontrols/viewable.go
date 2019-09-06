package xontrols

import "github.com/louisevanderlith/droxolite/context"

//Viewable handles controls that handle and view items.
type Viewable interface {
	Nomad
	View(context.Requester) (int, interface{})
}
