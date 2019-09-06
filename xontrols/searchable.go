package xontrols

import "github.com/louisevanderlith/droxolite/context"

//Searchable handles controls that handle search submissions[GET]
type Searchable interface {
	Nomad
	Search(context.Requester) (int, interface{})
}
