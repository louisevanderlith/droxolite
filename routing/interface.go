package routing

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

//NewInterfaceBundle returns a new RouteGroup that has been setup for UI purposes. The first will be used as the default
func NewInterfaceBundle(name string, required roletype.Enum, ctrls ...xontrols.InterfaceXontroller) *RouteGroup {
	if len(ctrls) == 0 {
		panic("ctrls must have at least one controller")
	}

	rg := NewRouteGroup(name, mix.Page)

	//Default Page
	deftName := fmt.Sprintf("Home %s", name)
	deftPath := fmt.Sprintf("/%s", strings.ToLower(name))
	rg.AddRoute(deftName, deftPath, http.MethodGet, required, ctrls[0].Default)

	for _, ctrl := range ctrls {
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

		if !searchable {
			rg.AddSubGroup(sub)
			continue
		}

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

		rg.AddSubGroup(sub)
	}

	return rg
}

func getControllerName(ctrl xontrols.InterfaceXontroller) string {
	tpe := reflect.TypeOf(ctrl).String()
	lstDot := strings.LastIndex(tpe, ".")

	if lstDot != -1 {
		return tpe[(lstDot + 1):]
	}

	return tpe
}
