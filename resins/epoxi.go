package resins

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/routing"
)

type Epoxi interface {
	Router() http.Handler
	Service() *bodies.Service
	AddGroup(routeGroup *routing.RouteGroup)
	Handle(mxFunc routing.MixerFunc, route *routing.Route) http.HandlerFunc
	EnableCORS(host string)
}
