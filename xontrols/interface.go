package xontrols

import "github.com/louisevanderlith/droxolite/context"

//Queries signals that the path can accept querystrings
type Queries interface {
	Nomad
	AcceptsQuery() map[string]string
}

//SearchableXontroller handles controls that handle search submissions and view items.
type Searchable interface {
	Nomad
	Search(context.Requester) (int, interface{})
	View(context.Requester) (int, interface{})
}

//CreateableXontroller handles controls that create content
type Createable interface {
	Searchable
	Create(context.Requester) (int, interface{})
}
