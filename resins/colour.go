package resins

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/element"
	"github.com/louisevanderlith/droxolite/filters"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

type ColourEpoxy struct {
	service  *bodies.Service
	router   http.Handler
	identity *element.Identity
	sideMenu *bodies.Menu
	//templates   *template.Template --rather use
	securityUrl string
	mxFunc      mix.InitFunc
}

//NewColourExpoxy returns a new Instance of the Epoxy with a Theme
func NewColourEpoxy(service *bodies.Service, d *element.Identity, securityUrl string, indexPage ServeFunc) Epoxi {

	e := &ColourEpoxy{
		service:     service,
		identity:    d,
		securityUrl: securityUrl,
		sideMenu:    bodies.NewMenu(),
		mxFunc:      mix.Page,
	}

	routr := mux.NewRouter()
	routr.HandleFunc("/", e.Handle("index", roletype.Unknown, indexPage))
	//Applications have assets in the 'dist' folder
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	routr.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	e.router = routr

	return e
}

/*
func (e *ColourEpoxy) AddBundle(b routing.Bundler) {
	interBun := b.(routing.InterfaceBundle)
	log.Println("RouteGroup:", b.RouteGroup().Name)

	sub := e.router.(*mux.Router).PathPrefix("/" + strings.ToLower(b.RouteGroup().Name)).Subrouter()
	e.sideMenu.AddGroup(b.RouteGroup().Name, interBun.SideMenu)

	for _, v := range b.RouteGroup().Routes {
		log.Println("Route:", v.Path)
		r := sub.Handle(v.Path, e.Handle(b.RouteGroup().MixFunc, v)).Methods(v.Method)

		for qkey, qval := range v.Queries {
			r.Queries(qkey, qval)
		}
	}

	//add sub groups
	for _, sgroup := range b.RouteGroup().SubGroups {
		log.Println("SubRoute:", strings.ToLower(sgroup.Name))
		xsub := sub.PathPrefix("/" + strings.ToLower(sgroup.Name)).Subrouter()

		for _, v := range sgroup.Routes {
			r := xsub.Handle(v.Path, e.Handle(b.RouteGroup().MixFunc, v)).Methods(v.Method)

			for qkey, qval := range v.Queries {
				r.Queries(qkey, qval)
			}
		}
	}
}*/

//JoinBundle will populate a mux.Router with paths generated from the ctrls names.
func (e *ColourEpoxy) JoinBundle(name string, required roletype.Enum, ctrls ...xontrols.Nomad) {
	if len(ctrls) == 0 {
		panic("ctrls must have at least one controller")
	}

	//The subrouter will create the basepath for every controller in the bundle
	//eg. Blog will create /blog
	subroute := e.router.(*mux.Router).PathPrefix("/" + name).Subrouter()
	subroute.Handle("", e.Handle("", roletype.Unknown, ctrls[0].Get)).Methods(http.MethodGet)

	var menu []bodies.MenuItem
	for _, ctrl := range ctrls {
		ctrlName := getControllerName(ctrl)
		ctrlPath := "/" + strings.ToLower(ctrlName)
		log.Println("Controller:", ctrlName)

		//The nested subrouter will create the basepath for every function in the controller
		//eg. Articles will create /blog/articles
		xsub := subroute.PathPrefix(ctrlPath).Subrouter()

		qryCtrl, isQueried := ctrl.(xontrols.Queries)

		if isQueried {
			for qkey, qval := range qryCtrl.AcceptsQuery() {
				xsub.Queries(qkey, qval)
			}
		}

		//Default
		xsub.Handle("", e.Handle(ctrlName, required, ctrl.Get)).Methods(http.MethodGet)

		searchCtrl, searchable := ctrl.(xontrols.Searchable)
		createable := false
		if searchable {
			//Search, it uses the default page
			xsub.Handle("/{pagesize:[A-Z][0-9]+}", e.Handle(ctrlName, required, searchCtrl.Search)).Methods(http.MethodGet)
			xsub.Handle("/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", e.Handle(ctrlName, required, searchCtrl.Search)).Methods(http.MethodGet)

			//View
			xsub.Handle("/{key:[0-9]+\x60[0-9]+}", e.Handle(ctrlName+"View", required, searchCtrl.View)).Methods(http.MethodGet)

			//Create
			createCtrl, createable := searchCtrl.(xontrols.Createable)

			if createable {
				xsub.Handle("/create", e.Handle(ctrlName+"Create", required, createCtrl.Create)).Methods(http.MethodGet)
			}
		}

		menu = append(menu, bodies.NewItem("", ctrlPath, ctrlName, createable, nil))
	}

	e.sideMenu.AddGroup(name, menu)

	/*var menu []bodies.MenuItem
	rg := NewRouteGroup(name, mix.Page)

	for _, ctrl := range ctrls {
		createable := false
		ctrlName := getControllerName(ctrl)
		sub := NewRouteGroup(ctrlName, mix.Page)

		//Default
		qryCtrl, isQueried := ctrl.(xontrols.QueriesXontrol)

		if isQueried {
			sub.AddRouteWithQueries("", "/", http.MethodGet, required, qryCtrl.AcceptsQuery(), ctrl.Default)
		} else {
			sub.AddRoute("", "/", http.MethodGet, required, ctrl.Default)
		}

		searchCtrl, searchable := ctrl.(xontrols.SearchableXontroller)

		if searchable {
			//Search, it uses the default page
			sub.AddRoute("", "/{pagesize:[A-Z][0-9]+}", http.MethodGet, required, searchCtrl.Search)
			sub.AddRoute("", "/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", http.MethodGet, required, searchCtrl.Search)

			//View
			sub.AddRoute("View", "/{key:[0-9]+\x60[0-9]+}", http.MethodGet, required, searchCtrl.View)

			//Create
			createCtrl, createable := searchCtrl.(xontrols.CreateableXontroller)

			if createable {
				sub.AddRoute("Create", "/create", http.MethodGet, required, createCtrl.Create)
			}
		}

		menu = append(menu, bodies.NewItem("", "/"+strings.ToLower(ctrlName), ctrlName, createable, nil))
		rg.AddSubGroup(sub)
	}

	return InterfaceBundle{menu, rg}*/
}

func getControllerName(ctrl xontrols.Nomad) string {
	tpe := reflect.TypeOf(ctrl).String()
	lstDot := strings.LastIndex(tpe, ".")

	if lstDot != -1 {
		return tpe[(lstDot + 1):]
	}

	return tpe
}

func (e *ColourEpoxy) Handle(name string, required roletype.Enum, process ServeFunc /*route *routing.Route*/) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, e.service.ID)

		allow, avoc := filters.TokenCookieCheck(ctx, required, e.service.PublicKey, e.service.Name)
		if !allow {
			err := sendToLogin(ctx, e.securityUrl)

			if err != nil {
				log.Panicln(err)
			}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := process(ctx)
		mxer := e.mxFunc(name, data, e.identity, avoc).(mix.ColourMixer)
		//mxer.ApplySettings(route.Name, *e.settings, avoc)

		mxer.CreateSideMenu(e.sideMenu)
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

func sendToLogin(ctx context.Contexer, securityURL string) error {
	scheme := ctx.Scheme()

	if len(scheme) == 0 {
		scheme = "https"
	}

	moveURL := fmt.Sprintf("%s://%s%s", scheme, ctx.Host(), ctx.RequestURI())
	loginURL := buildLoginURL(securityURL, moveURL)

	ctx.Redirect(http.StatusTemporaryRedirect, loginURL)

	return nil
}

func buildLoginURL(securityURL, returnURL string) string {
	cleanReturn := removeQueries(returnURL)
	escURL := url.QueryEscape(cleanReturn)
	return fmt.Sprintf("%slogin?return=%s", securityURL, escURL)
}

func removeQueries(url string) string {
	idxOfQuery := strings.Index(url, "?")

	if idxOfQuery != -1 {
		url = url[:idxOfQuery]
	}

	return url
}

func buildSubscribeURL(securityURL string) string {
	return fmt.Sprintf("%ssubscribe", securityURL)
}
