package resins

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/filters"
	"github.com/louisevanderlith/droxolite/routing"
)

type ColourEpoxy struct {
	service     *bodies.Service
	router      http.Handler
	settings    *bodies.ThemeSetting
	sideMenu    *bodies.Menu
	masterpage  string
	templates   *template.Template
	securityUrl string
}

//NewColourExpoxy returns a new Instance of the Epoxy with a Theme
func NewColourEpoxy(service *bodies.Service, settings bodies.ThemeSetting, masterpage, securityUrl string) Epoxi {
	routr := mux.NewRouter()

	//Applications have assets in the 'dist' folder
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	routr.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	e := &ColourEpoxy{
		service:    service,
		router:     routr,
		settings:   &settings,
		masterpage: masterpage,
		sideMenu:   bodies.NewMenu(),
	}

	err := e.settings.LoadTemplate("./views", masterpage)

	if err != nil {
		panic(err)
	}

	return e
}

func (e *ColourEpoxy) AddGroup(routeGroup *routing.RouteGroup) {
	if e.settings == nil {
		log.Fatalf("Use the Colour Epoxy!")
	}

	var menuGroup []bodies.MenuItem
	for _, v := range routeGroup.Routes {
		if v.Method == http.MethodGet {
			//menuGroup = append(menuGroup, bodies.NewItem(, v.Path, v.Name, nil))
			//menuGroup.AddItem(v.Path, v.Name, nil)
			baseURL := v.Path

			for k, sgroup := range routeGroup.SubGroups {
				subPath := baseURL + "/" + strings.ToLower(sgroup.Name)
				var subChildren []bodies.MenuItem

				for sk, sv := range sgroup.Routes {
					if sv.Method == http.MethodGet && !strings.HasPrefix(sv.Path, "/{") && sv.Name != "Default" {
						subChildren = append(subChildren, bodies.NewItem(fmt.Sprintf("c%v", sk), subPath+sv.Path, sv.Name, nil))
					}
				}

				subMenu := bodies.NewItem(fmt.Sprintf("r%v", k), subPath, sgroup.Name, subChildren)
				menuGroup = append(menuGroup, subMenu)
			}
		}
	}

	e.sideMenu.AddGroup(routeGroup.Name, menuGroup)

	sub := e.router.(*mux.Router).PathPrefix("/" + strings.ToLower(routeGroup.Name)).Subrouter()

	for _, v := range routeGroup.Routes {
		r := sub.Handle(v.Path, e.Handle(routeGroup.MixFunc, v)).Methods(v.Method)

		for qkey, qval := range v.Queries {
			r.Queries(qkey, qval)
		}
	}

	//add sub groups
	for _, sgroup := range routeGroup.SubGroups {
		xsub := sub.PathPrefix("/" + strings.ToLower(sgroup.Name)).Subrouter()

		for _, v := range sgroup.Routes {
			r := xsub.Handle(v.Path, e.Handle(routeGroup.MixFunc, v)).Methods(v.Method)

			for qkey, qval := range v.Queries {
				r.Queries(qkey, qval)
			}
		}
	}
}

func (e *ColourEpoxy) Handle(mxFunc routing.MixerFunc, route *routing.Route) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, e.service.ID)

		if !filters.TokenCookieCheck(ctx, route.RequiredRole, e.service.PublicKey, e.service.Name) {
			//err := sendToLogin(ctrl.Ctx(), e.securityUrl)

			//if err != nil {
			//	log.Panicln(err)
			//}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := route.Function(ctx)
		mxer := mxFunc(data)

		mxer.ApplySettings(route.Name, *e.settings)
		err := ctx.Serve(status, mxer)

		if err != nil {
			log.Panicln(err)
		}
	}
}

func (e *ColourEpoxy) Router() http.Handler {
	return e.router
}

func (e *ColourEpoxy) Service() *bodies.Service {
	return e.service
}

func (e *ColourEpoxy) EnableCORS(host string) {
	//No Need.
}
