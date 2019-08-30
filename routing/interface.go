package routing

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

type InterfaceBundle struct {
	SideMenu   []bodies.MenuItem
	routeGroup *RouteGroup
}

func (i InterfaceBundle) RouteGroup() *RouteGroup {
	return i.routeGroup
}

//NewInterfaceBundle returns a new RouteGroup that has been setup for UI purposes. The first will be used as the default
func NewInterfaceBundle(name string, required roletype.Enum, ctrls ...xontrols.InterfaceXontroller) InterfaceBundle {
	if len(ctrls) == 0 {
		panic("ctrls must have at least one controller")
	}

	result := InterfaceBundle{
		routeGroup: NewRouteGroup(name, mix.Page),
	}

	//Default Page
	//deftName := fmt.Sprintf("Home %s", name)
	//deftPath := fmt.Sprintf("/%s", strings.ToLower(name))
	//result.routeGroup.AddRoute(deftName, deftPath, http.MethodGet, required, ctrls[0].Default)
	//var menuGroup []bodies.MenuItem
	for _, ctrl := range ctrls {
		createable := false
		ctrlName := getControllerName(ctrl)

		sub := NewRouteGroup(ctrlName, mix.Page)

		//Default
		qryCtrl, isQueried := ctrl.(xontrols.QueriesXontrol)

		if isQueried {
			sub.AddRouteWithQueries(ctrlName, "", http.MethodGet, required, qryCtrl.AcceptsQuery(), ctrl.Default)
		} else {
			sub.AddRoute(ctrlName, "", http.MethodGet, required, ctrl.Default)
		}

		searchCtrl, searchable := ctrl.(xontrols.SearchableXontroller)

		if searchable {
			//Search, it uses the default page
			sub.AddRoute(ctrlName, "/{pagesize:[A-Z][0-9]+}", http.MethodGet, required, searchCtrl.Search)
			sub.AddRoute(ctrlName, "/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", http.MethodGet, required, searchCtrl.Search)

			//View
			sub.AddRoute(ctrlName+"View", "/{key:[0-9]+\x60[0-9]+}", http.MethodGet, required, searchCtrl.View)

			//Create
			createCtrl, createable := searchCtrl.(xontrols.CreateableXontroller)

			if createable {
				sub.AddRoute(ctrlName+"Create", "/create", http.MethodGet, required, createCtrl.Create)
			}
		}

		result.SideMenu = append(result.SideMenu, bodies.NewItem("", "/"+strings.ToLower(ctrlName), ctrlName, createable, nil))

		result.routeGroup.AddSubGroup(sub)
	}

	return result
}

func getControllerName(ctrl xontrols.InterfaceXontroller) string {
	tpe := reflect.TypeOf(ctrl).String()
	lstDot := strings.LastIndex(tpe, ".")

	if lstDot != -1 {
		return tpe[(lstDot + 1):]
	}

	return tpe
}
