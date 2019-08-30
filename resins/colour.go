package resins

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/element"
	"github.com/louisevanderlith/droxolite/filters"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/routing"
)

type ColourEpoxy struct {
	service  *bodies.Service
	router   http.Handler
	identity *element.Identity
	sideMenu *bodies.Menu
	//templates   *template.Template --rather use
	securityUrl string
}

//NewColourExpoxy returns a new Instance of the Epoxy with a Theme
func NewColourEpoxy(service *bodies.Service, d *element.Identity, securityUrl string) Epoxi {
	routr := mux.NewRouter()

	//Applications have assets in the 'dist' folder
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	routr.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	e := &ColourEpoxy{
		service:     service,
		router:      routr,
		identity:    d,
		securityUrl: securityUrl,
		sideMenu:    bodies.NewMenu(),
	}

	return e
}

func (e *ColourEpoxy) AddBundle(b routing.Bundler) {
	interBun := b.(routing.InterfaceBundle)

	sub := e.router.(*mux.Router).PathPrefix("/" + strings.ToLower(b.RouteGroup().Name)).Subrouter()
	e.sideMenu.AddGroup(b.RouteGroup().Name, interBun.SideMenu)
	for _, v := range b.RouteGroup().Routes {
		r := sub.Handle(v.Path, e.Handle(b.RouteGroup().MixFunc, v)).Methods(v.Method)

		for qkey, qval := range v.Queries {
			r.Queries(qkey, qval)
		}
	}

	//add sub groups
	for _, sgroup := range b.RouteGroup().SubGroups {
		xsub := sub.PathPrefix("/" + strings.ToLower(sgroup.Name)).Subrouter()

		for _, v := range sgroup.Routes {
			r := xsub.Handle(v.Path, e.Handle(b.RouteGroup().MixFunc, v)).Methods(v.Method)

			for qkey, qval := range v.Queries {
				r.Queries(qkey, qval)
			}
		}
	}
}

func (e *ColourEpoxy) Handle(mxFunc mix.InitFunc, route *routing.Route) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, e.service.ID)

		allow, avoc := filters.TokenCookieCheck(ctx, route.RequiredRole, e.service.PublicKey, e.service.Name)
		if !allow {
			err := sendToLogin(ctx, e.securityUrl)

			if err != nil {
				log.Panicln(err)
			}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := route.Function(ctx)
		mxer := mxFunc(route.Name, data, e.identity, avoc).(mix.ColourMixer)
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
