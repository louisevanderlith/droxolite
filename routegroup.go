package droxolite

import (
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

type Route struct {
	Name         string //Name will be used to display links.
	Path         string
	Method       string
	RequiredRole roletype.Enum
	Queries      map[string]string
	Function     func()
}

type RouteGroup struct {
	Name       string
	Controller xontrols.Controller
	Routes     []*Route
	SubGroups  []*RouteGroup
}

func NewRouteGroup(name string, ctrl xontrols.Controller) *RouteGroup {
	return &RouteGroup{
		Name:       name,
		Controller: ctrl,
	}
}

func (g *RouteGroup) AddSubGroup(subgrp *RouteGroup) {
	g.SubGroups = append(g.SubGroups, subgrp)
}

func (g *RouteGroup) AddRoute(name, path, method string, requiredRole roletype.Enum, function func()) *Route {
	result := &Route{
		Name:         name,
		Path:         path,
		Method:       method,
		RequiredRole: requiredRole,
		Function:     function,
		Queries:      make(map[string]string),
	}

	g.Routes = append(g.Routes, result)

	return result
}

func (g *RouteGroup) AddRouteWithQueries(name, path, method string, requiredRole roletype.Enum, queries map[string]string, function func()) *Route {
	result := &Route{
		Name:         name,
		Path:         path,
		Method:       method,
		RequiredRole: requiredRole,
		Function:     function,
		Queries:      queries,
	}

	g.Routes = append(g.Routes, result)

	return result
}
