package resins

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/element"
	"github.com/louisevanderlith/droxolite/filters"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/routing"
	"github.com/rs/cors"
)

type BasicEpoxy struct {
	service  *bodies.Service
	router   http.Handler
	identity *element.Identity
}

//NewBasicExpoxy returns a new Instance of the Epoxy
func NewBasicEpoxy(service *bodies.Service, d *element.Identity) Epoxi {
	routr := mux.NewRouter()

	return &BasicEpoxy{
		service:  service,
		router:   routr,
		identity: d,
	}
}

func (e *BasicEpoxy) AddBundle(b routing.Bundler) {
	sub := e.router.(*mux.Router).PathPrefix("/" + strings.ToLower(b.RouteGroup().Name)).Subrouter()

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

func (e *BasicEpoxy) Handle(mxFunc mix.InitFunc, route *routing.Route) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, e.service.ID)

		allow, avoc := filters.TokenCheck(ctx, route.RequiredRole, e.service.PublicKey, e.service.Name)
		if !allow {
			err := ctx.Serve(http.StatusUnauthorized, mxFunc(route.Name, nil, e.identity, nil))

			if err != nil {
				log.Panicln(err)
			}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := route.Function(ctx)
		mxer := mxFunc(ctx.RequestURI(), data, e.identity, avoc)

		//mxer.ApplySettings(ctx.RequestURI(), *e.settings, avoc)
		err := ctx.Serve(status, mxer)
		if err != nil {
			log.Panicln(err)
		}
	}
}

func (e *BasicEpoxy) Router() http.Handler {
	return e.router
}

func (e *BasicEpoxy) Service() *bodies.Service {
	return e.service
}

func (e *BasicEpoxy) EnableCORS(host string) {
	allowed := fmt.Sprintf("https://*%s", strings.TrimSuffix(host, "/"))

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{allowed}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application
		},
	})

	e.router = corsOpts.Handler(e.router)
}
