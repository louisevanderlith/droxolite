package resins

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

type Epoxi interface {
	Router() http.Handler
	Service() *bodies.Service
	JoinXontrol(path string, required roletype.Enum, mxFunc mix.InitFunc, ctrl xontrols.Nomad)
	JoinBundle(path string, required roletype.Enum, mxFunc mix.InitFunc, ctrls ...xontrols.Nomad)
	EnableCORS(host string)
}

type ServeFunc func(context.Requester) (int, interface{})
