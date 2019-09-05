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
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
	"github.com/rs/cors"
)

type BasicEpoxy struct {
	service  *bodies.Service
	router   http.Handler
	identity *element.Identity
	mxFunc   mix.InitFunc
}

//NewBasicExpoxy returns a new Instance of the Epoxy
func NewBasicEpoxy(service *bodies.Service, d *element.Identity, mxFunc mix.InitFunc) Epoxi {
	routr := mux.NewRouter()

	return &BasicEpoxy{
		service:  service,
		router:   routr,
		identity: d,
		mxFunc:   mxFunc,
	}
}

func (e *BasicEpoxy) JoinBundle(name string, required roletype.Enum, ctrls ...xontrols.Nomad) {
	sub := e.router.(*mux.Router).PathPrefix("/" + strings.ToLower(name)).Subrouter()

	for _, ctrl := range ctrls {
		ctrlName := getControllerName(ctrl)
		ctrlPath := "/" + strings.ToLower(ctrlName)
		log.Println("Controller:", ctrlName)

		//The nested subrouter will create the basepath for every function in the controller
		//eg. Articles will create /blog/articles
		xsub := sub.PathPrefix(ctrlPath).Subrouter()

		//Default
		xsub.Handle("", e.Handle(ctrlName, required, ctrl.Get)).Methods(http.MethodGet)

		//Storable
		storeCtrl, isStore := ctrl.(xontrols.Store)

		if isStore {
			xsub.Handle("", e.Handle(ctrlName+"Create", required, storeCtrl.Create)).Methods(http.MethodPost)

			xsub.Handle("/{key:[0-9]+\x60[0-9]+}", e.Handle(ctrlName, required, storeCtrl.GetOne)).Methods(http.MethodGet)
			xsub.Handle("/{key:[0-9]+\x60[0-9]+}", e.Handle(ctrlName, required, storeCtrl.Update)).Methods(http.MethodPut)
			xsub.Handle("/{key:[0-9]+\x60[0-9]+}", e.Handle(ctrlName, required, storeCtrl.Delete)).Methods(http.MethodDelete)

			xsub.Handle("/all/{pagesize:[A-Z][0-9]+}", e.Handle("Get All", required, storeCtrl.Get)).Methods(http.MethodGet)

			qryCtrl, isQueried := ctrl.(xontrols.Queries)

			if isQueried {
				for qkey, qval := range qryCtrl.AcceptsQuery() {
					xsub.Queries(qkey, qval)
				}
			}
		}
	}
}

func (e *BasicEpoxy) Handle(name string, required roletype.Enum, process ServeFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, e.service.ID)

		allow, avoc := filters.TokenCheck(ctx, required, e.service.PublicKey, e.service.Name)
		if !allow {
			err := ctx.Serve(http.StatusUnauthorized, e.mxFunc(name, nil, e.identity, nil))

			if err != nil {
				log.Panicln(err)
			}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := process(ctx)
		mxer := e.mxFunc(ctx.RequestURI(), data, e.identity, avoc)

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
