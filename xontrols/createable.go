package xontrols

import "github.com/louisevanderlith/droxolite/context"

//Createable handles controls that create content [POST]
type Createable interface {
	Searchable
	Create(context.Requester) (int, interface{})
}
