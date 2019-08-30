package resins

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/routing"
)

type Epoxi interface {
	Router() http.Handler
	Service() *bodies.Service
	AddBundle(b routing.Bundler)
	Handle(mxFunc mix.InitFunc, route *routing.Route) http.HandlerFunc
	EnableCORS(host string)
}
