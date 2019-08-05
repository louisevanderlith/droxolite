package routing

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/louisevanderlith/droxolite"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

//NewInterfaceBundle returns a new RouteGroup that has been setup for UI purposes. The first will be used as the default
func NewInterfaceBundle(name string, required roletype.Enum, ctrls ...xontrols.InterfaceXontroller) *droxolite.RouteGroup {
	if len(ctrls) == 0 {
		panic("ctrls must have at least one controller")
	}

	rg := droxolite.NewRouteGroup(name, ctrls[0].(xontrols.Controller))

	//Default Page
	deftName := fmt.Sprintf("Default %s", name)
	deftPath := fmt.Sprintf("/%s", strings.ToLower(name))
	rg.AddRoute(deftName, deftPath, http.MethodGet, required, ctrls[0].Default)

	for _, ctrl := range ctrls {
		ctrlName := getControllerName(ctrl)
		sub := droxolite.NewRouteGroup(ctrlName, ctrl.(xontrols.Controller))

		searchCtrl, searchable := ctrl.(xontrols.SearchableXontroller)

		if !searchable {
			continue
		}

		//Search
		sub.AddRoute("Search", "/{pagesize:[A-Z][0-9]+}", http.MethodGet, required, searchCtrl.Search)
		sub.AddRoute("Search Query", "/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", http.MethodGet, required, searchCtrl.Search)

		//View
		sub.AddRoute("View", "/{key:[0-9]+\x60[0-9]+}", http.MethodGet, required, searchCtrl.View)

		//Create
		createCtrl, createable := searchCtrl.(xontrols.CreateableXontroller)

		if createable {
			sub.AddRoute("Create", "/create", http.MethodGet, required, createCtrl.Create)
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

/*

	//Stock
	stockGroup := droxolite.NewRouteGroup("Stock", nil)
	carsCtrl := &stock.CarsController{}
	carsGroup := droxolite.NewRouteGroup("Cars", carsCtrl)
	carsGroup.AddRoute("All Cars", "/{pagesize:[A-Z][0-9]+}", "GET", roletype.Admin, carsCtrl.Get)
	carsGroup.AddRoute("Edit Car", "/edit/{key:[0-9]+\x60[0-9]+}", "GET", roletype.Admin, carsCtrl.GetEdit)
	stockGroup.AddSubGroup(carsGroup)

	partsCtrl := &stock.PartsController{}
	partsGroup := droxolite.NewRouteGroup("Parts", partsCtrl)
	partsGroup.AddRoute("All Parts", "/{pagesize:[A-Z][0-9]+}", "GET", roletype.Admin, partsCtrl.Get)
	partsGroup.AddRoute("Edit Part", "/edit/{key:[0-9]+\x60[0-9]+}", "GET", roletype.Admin, partsCtrl.GetEdit)
	stockGroup.AddSubGroup(partsGroup)

	srvcCtrl := &stock.ServicesController{}
	srvcGroup := droxolite.NewRouteGroup("Services", srvcCtrl)
	srvcGroup.AddRoute("All Services", "/{pagesize:[A-Z][0-9]+}", "GET", roletype.Admin, srvcCtrl.Get)
	srvcGroup.AddRoute("Edit Service", "/edit/{key:[0-9]+\x60[0-9]+}", "GET", roletype.Admin, srvcCtrl.GetEdit)
	stockGroup.AddSubGroup(srvcGroup)

	e.AddNamedGroup("Stock.API", stockGroup)
*/
