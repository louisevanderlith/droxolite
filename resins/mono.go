package resins

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/element"
	"github.com/louisevanderlith/droxolite/filters"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
	"github.com/louisevanderlith/proofclient/models"
	"github.com/rs/cors"
)

type monoEpoxy struct {
	clientCred models.ClientCred
	router     http.Handler
	identity   *element.Identity
}

func NewMonoEpoxy(clientCred models.ClientCred, d *element.Identity) Epoxi {
	routr := mux.NewRouter()

	return &monoEpoxy{
		client:   client,
		router:   routr,
		identity: d,
	}
}

//Routers returns the final http Handler which can be listened on
func (e *monoEpoxy) Router() http.Handler {
	return e.router
}

func (e *monoEpoxy) Client() models.ClientCred {
	return e.client
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

func (e *monoEpoxy) JoinPath(r *mux.Router, path, name, method string, required roletype.Enum, mxFunc mix.InitFunc, f ServeFunc) {
	r.Handle(path, e.filter(name, required, mxFunc, f)).Methods(method)
}

func (e *monoEpoxy) JoinBundle(path string, required roletype.Enum, mxFunc mix.InitFunc, ctrls ...xontrols.Nomad) {
	sub := e.router.(*mux.Router).PathPrefix(path).Subrouter()

	for _, ctrl := range ctrls {
		ctrlName := getControllerName(ctrl)
		ctrlPath := "/" + strings.ToLower(ctrlName)
		ctrlSub := sub.PathPrefix(ctrlPath).Subrouter()

		//Get
		e.JoinPath(ctrlSub, "", ctrlName, http.MethodGet, required, mxFunc, ctrl.Get)

		//Search
		searchCtrl, isSearch := ctrl.(xontrols.Searchable)

		if isSearch {
			e.JoinPath(ctrlSub, "/{pagesize:[A-Z][0-9]+}", ctrlName, http.MethodGet, required, mxFunc, searchCtrl.Search)
			e.JoinPath(ctrlSub, "/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", ctrlName, http.MethodGet, required, mxFunc, searchCtrl.Search)
		}

		//View
		viewCtrl, isView := ctrl.(xontrols.Viewable)

		if isView {
			e.JoinPath(ctrlSub, "/{key:[0-9]+\x60[0-9]+}", ctrlName, http.MethodGet, required, mxFunc, viewCtrl.View)
		}

		//Create
		createCtrl, isCreate := ctrl.(xontrols.Createable)

		if isCreate {
			e.JoinPath(ctrlSub, "", ctrlName, http.MethodPost, required, mxFunc, createCtrl.Create)
		}

		//Update
		updatCtrl, isUpdate := ctrl.(xontrols.Updateable)

		if isUpdate {
			e.JoinPath(ctrlSub, "/{key:[0-9]+\x60[0-9]+}", ctrlName, http.MethodPut, required, mxFunc, updatCtrl.Update)
		}

		//Delete
		delCtrl, isDelete := ctrl.(xontrols.Deleteable)

		if isDelete {
			e.JoinPath(ctrlSub, "/{key:[0-9]+\x60[0-9]+}", ctrlName, http.MethodDelete, required, mxFunc, delCtrl.Delete)
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
		ctx := context.New(resp, req, srv.ID, srv.PublicKey)
		
		p := filters.Pack{
			RequestURI: ctx.
		}
		/*
		RequestURI   string
	Token        string
	RequiredRole roletype.Enum
	ClientName   string //used to be service name
	ClientCred   models.ClientCred
	Inspector    client.Inspector
		*/
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
