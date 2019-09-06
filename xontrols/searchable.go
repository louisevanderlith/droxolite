package xontrols

import "github.com/louisevanderlith/droxolite/context"

//Searchable handles controls that handle search submissions and view items. [GET]
type Searchable interface {
	Nomad
	Search(context.Requester) (int, interface{})
	View(context.Requester) (int, interface{})
}
