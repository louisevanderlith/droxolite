package xontrols

//InterfaceXontroller can only handle GET Requests
//The interface should be used as a guide when creating page groups
type InterfaceXontroller interface {
	UIController
	Default() error
}

//SearchableXontroller handles controls that handle search submissions and view items.
type SearchableXontroller interface {
	InterfaceXontroller
	Search() error
	View() error
}

//CreateableXontroller handles controls that create content
type CreateableXontroller interface {
	SearchableXontroller
	Create() error
}
