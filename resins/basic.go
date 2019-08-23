package resins

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/filters"
	"github.com/louisevanderlith/droxolite/routing"
)

type BasicEpoxy struct {
	service  *bodies.Service
	router   http.Handler
	settings interface{}
}

//NewBasicExpoxy returns a new Instance of the Epoxy
func NewBasicEpoxy(service *bodies.Service) Epoxi {
	routr := mux.NewRouter()

	return &BasicEpoxy{
		service: service,
		router:  routr,
	}
}

func (e *BasicEpoxy) AddGroup(routeGroup *routing.RouteGroup) {
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

func (e *BasicEpoxy) Handle(mxFunc routing.MixerFunc, route *routing.Route) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, e.service.ID)

		if !filters.TokenCheck(ctx, route.RequiredRole, e.service.PublicKey, e.service.Name) {
			//err := sendToLogin(ctx, e.service.)

			//if err != nil {
			//	log.Panicln(err)
			//}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := route.Function(ctx)
		err := ctx.Serve(status, mxFunc(data))

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
