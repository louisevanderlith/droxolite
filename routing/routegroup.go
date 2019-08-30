package routing

import (
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/roletype"
)

type Route struct {
	Name         string //Name will be used to display links.
	Path         string
	PageName     string
	Method       string
	RequiredRole roletype.Enum
	Queries      map[string]string
	Function     ServeFunc
}

type RouteGroup struct {
	Name      string
	MixFunc   mix.InitFunc
	Routes    []*Route
	SubGroups []*RouteGroup
}

type ServeFunc func(context.Requester) (int, interface{})

func NewRouteGroup(name string, mxFunc mix.InitFunc) *RouteGroup {
	return &RouteGroup{
		Name:    name,
		MixFunc: mxFunc,
	}
}

func (g *RouteGroup) RouteGroup() *RouteGroup {
	return g
}

func (g *RouteGroup) AddSubGroup(subgrp *RouteGroup) {
	g.SubGroups = append(g.SubGroups, subgrp)
}

func (g *RouteGroup) AddRoute(name, path, method string, requiredRole roletype.Enum, function ServeFunc) *Route {
	result := &Route{
		Name:         name,
		PageName:     g.Name + name,
		Path:         path,
		Method:       method,
		RequiredRole: requiredRole,
		Function:     function,
		Queries:      make(map[string]string),
	}

	g.Routes = append(g.Routes, result)

	return result
}

func (g *RouteGroup) AddRouteWithQueries(name, path, method string, requiredRole roletype.Enum, queries map[string]string, function ServeFunc) *Route {
	result := &Route{
		Name:         name,
		PageName:     g.Name + name,
		Path:         path,
		Method:       method,
		RequiredRole: requiredRole,
		Function:     function,
		Queries:      queries,
	}

	g.Routes = append(g.Routes, result)

	return result
}
