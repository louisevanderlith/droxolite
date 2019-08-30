package xontrols

import "github.com/louisevanderlith/droxolite/context"

//InterfaceXontroller can only handle GET Requests
//The interface should be used as a guide when creating page groups
type InterfaceXontroller interface {
	Default(context.Requester) (int, interface{})
}

type QueriesXontrol interface {
	InterfaceXontroller
	AcceptsQuery() map[string]string
}

//SearchableXontroller handles controls that handle search submissions and view items.
type SearchableXontroller interface {
	InterfaceXontroller
	Search(context.Requester) (int, interface{})
	View(context.Requester) (int, interface{})
}

//CreateableXontroller handles controls that create content
type CreateableXontroller interface {
	SearchableXontroller
	Create(context.Requester) (int, interface{})
}
