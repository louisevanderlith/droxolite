package xontrols

//InterfaceXontroller can only handle GET Requests
//The interface should be used as a guide when creating page groups
type InterfaceXontroller interface {
	Default()
}

//SearchableXontroller handles controls that handle search submissions and view items.
type SearchableXontroller interface {
	InterfaceXontroller
	Search()
	View()
}

//CreateableXontroller handles controls that create content
type CreateableXontroller interface {
	SearchableXontroller
	Create()
}
