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

type monoEpoxy struct {
	service  *bodies.Service
	router   http.Handler
	identity *element.Identity
}

func NewMonoEpoxy(service *bodies.Service, d *element.Identity) Epoxi {
	routr := mux.NewRouter()

	return &monoEpoxy{
		service:  service,
		router:   routr,
		identity: d,
	}
}

//Routers returns the final http Handler which can be listened on
func (e *monoEpoxy) Router() http.Handler {
	return e.router
}

func (e *monoEpoxy) Service() *bodies.Service {
	return e.service
}

func (e *monoEpoxy) EnableCORS(host string) {
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

func (e *monoEpoxy) JoinXontrol(path string, required roletype.Enum, mxFunc mix.InitFunc, ctrl xontrols.Nomad) {
	sub := e.router.(*mux.Router).PathPrefix(path).Subrouter()

	ctrlName := getControllerName(ctrl)
	ctrlPath := "/" + strings.ToLower(ctrlName)
	ctrlSub := sub.PathPrefix(ctrlPath).Subrouter()

	//Get
	ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, ctrl.Get)).Methods(http.MethodGet)

	//Search & View
	searchCtrl, isSearch := ctrl.(xontrols.Searchable)

	if isSearch {
		ctrlSub.Handle("/{pagesize:[A-Z][0-9]+}", e.filter(ctrlName, required, mxFunc, searchCtrl.Search)).Methods(http.MethodGet)
		ctrlSub.Handle("/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", e.filter(ctrlName, required, mxFunc, searchCtrl.Search)).Methods(http.MethodGet)
		ctrlSub.Handle("/{key:[0-9]+\x60[0-9]+}", e.filter(ctrlName, required, mxFunc, searchCtrl.View)).Methods(http.MethodGet)
	}

	//Create
	createCtrl, isCreate := ctrl.(xontrols.Createable)

	if isCreate {
		ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, createCtrl.Create)).Methods(http.MethodPost)
	}

	//Update
	updatCtrl, isUpdate := ctrl.(xontrols.Updateable)

	if isUpdate {
		ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, updatCtrl.Update)).Methods(http.MethodPut)
	}

	//Delete
	delCtrl, isDelete := ctrl.(xontrols.Deleteable)

	if isDelete {
		ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, delCtrl.Delete)).Methods(http.MethodDelete)
	}

	//Queries
	qryCtrl, isQueried := ctrl.(xontrols.Queries)

	if isQueried {
		for qkey, qval := range qryCtrl.AcceptsQuery() {
			ctrlSub.Queries(qkey, qval)
		}
	}
}

func (e *monoEpoxy) JoinBundle(path string, required roletype.Enum, mxFunc mix.InitFunc, ctrls ...xontrols.Nomad) {
	sub := e.router.(*mux.Router).PathPrefix(path).Subrouter()

	for _, ctrl := range ctrls {
		ctrlName := getControllerName(ctrl)
		ctrlPath := "/" + strings.ToLower(ctrlName)
		ctrlSub := sub.PathPrefix(ctrlPath).Subrouter()

		//Get
		ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, ctrl.Get)).Methods(http.MethodGet)

		//Search & View
		searchCtrl, isSearch := ctrl.(xontrols.Searchable)

		if isSearch {
			ctrlSub.Handle("/{pagesize:[A-Z][0-9]+}", e.filter(ctrlName, required, mxFunc, searchCtrl.Search)).Methods(http.MethodGet)
			ctrlSub.Handle("/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", e.filter(ctrlName, required, mxFunc, searchCtrl.Search)).Methods(http.MethodGet)
			ctrlSub.Handle("/{key:[0-9]+\x60[0-9]+}", e.filter(ctrlName, required, mxFunc, searchCtrl.View)).Methods(http.MethodGet)
		}

		//Create
		createCtrl, isCreate := ctrl.(xontrols.Createable)

		if isCreate {
			ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, createCtrl.Create)).Methods(http.MethodPost)
		}

		//Update
		updatCtrl, isUpdate := ctrl.(xontrols.Updateable)

		if isUpdate {
			ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, updatCtrl.Update)).Methods(http.MethodPut)
		}

		//Delete
		delCtrl, isDelete := ctrl.(xontrols.Deleteable)

		if isDelete {
			ctrlSub.Handle("", e.filter(ctrlName, required, mxFunc, delCtrl.Delete)).Methods(http.MethodDelete)
		}

		//Queries
		qryCtrl, isQueried := ctrl.(xontrols.Queries)

		if isQueried {
			for qkey, qval := range qryCtrl.AcceptsQuery() {
				ctrlSub.Queries(qkey, qval)
			}
		}
	}
}

func (e *monoEpoxy) filter(name string, required roletype.Enum, mxFunc mix.InitFunc, process ServeFunc) http.HandlerFunc {
	srv := e.service
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, srv.ID)

		allow, avoc := filters.TokenCheck(ctx, required, srv.PublicKey, srv.Name)
		if !allow {
			err := ctx.Serve(http.StatusUnauthorized, mxFunc(name, nil, e.identity, nil))

			if err != nil {
				log.Panicln(err)
			}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := process(ctx)
		mxer := mxFunc(ctx.RequestURI(), data, e.identity, avoc)

		err := ctx.Serve(status, mxer)
		if err != nil {
			log.Panicln(err)
		}
	}
}
