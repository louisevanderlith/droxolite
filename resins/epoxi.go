package resins

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

type Epoxi interface {
	Router() http.Handler
	Service() *bodies.Service
	JoinBundle(name string, required roletype.Enum, ctrls ...xontrols.Nomad)
	Handle(name string, required roletype.Enum, process ServeFunc) http.HandlerFunc
	EnableCORS(host string)
}

type ServeFunc func(context.Requester) (int, interface{})
