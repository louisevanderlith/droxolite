package resins

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/element"
	"github.com/louisevanderlith/droxolite/filters"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

type colourEpoxy struct {
	service     *bodies.Service
	router      http.Handler
	identity    *element.Identity
	sideMenu    *bodies.Menu
	securityUrl string
}

//NewColourExpoxy returns a new Instance of the Epoxy with a Theme
func NewColourEpoxy(service *bodies.Service, d *element.Identity, securityUrl string, indexPage ServeFunc) Epoxi {

	e := &colourEpoxy{
		service:     service,
		identity:    d,
		securityUrl: securityUrl,
		sideMenu:    bodies.NewMenu(),
	}

	routr := mux.NewRouter()
	routr.HandleFunc("/", e.filter("index", roletype.Unknown, mix.Page, indexPage))
	//Applications have assets in the 'dist' folder
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	routr.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	e.router = routr

	return e
}

func (e *colourEpoxy) Router() http.Handler {
	return e.router
}

func (e *colourEpoxy) Service() *bodies.Service {
	return e.service
}

func (e *colourEpoxy) EnableCORS(host string) {
	//No Need.
}

func (e *colourEpoxy) JoinXontrol(path string, required roletype.Enum, mxFunc mix.InitFunc, ctrl xontrols.Nomad) {
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

	//Update - Pages can't update
	//Delete -Pages can't delete

	//Queries
	qryCtrl, isQueried := ctrl.(xontrols.Queries)

	if isQueried {
		for qkey, qval := range qryCtrl.AcceptsQuery() {
			ctrlSub.Queries(qkey, qval)
		}
	}

	menuPath := ctrlPath
	groupName := "General"

	//more than just a slash
	if strings.HasPrefix(path, "/") && len(path) > 1 {
		menuPath = path + ctrlPath
		groupName = strings.Title(strings.Replace(path, "/", "", -1))
	}

	e.sideMenu.AddGroup(groupName, []bodies.MenuItem{bodies.NewItem("", menuPath, ctrlName, isCreate, nil)})
}

func (e *colourEpoxy) JoinBundle(path string, required roletype.Enum, mxFunc mix.InitFunc, ctrls ...xontrols.Nomad) {
	sub := e.router.(*mux.Router).PathPrefix(path).Subrouter()
	var menu []bodies.MenuItem

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
			ctrlSub.Handle("/{key:[0-9]+\x60[0-9]+}", e.filter(ctrlName+"View", required, mxFunc, searchCtrl.View)).Methods(http.MethodGet)
		}

		//Create
		createCtrl, isCreate := ctrl.(xontrols.Createable)

		if isCreate {
			ctrlSub.Handle("/create", e.filter(ctrlName, required, mxFunc, createCtrl.Create)).Methods(http.MethodGet)
		}

		//Update - Pages can't update
		//Delete -Pages can't delete

		//Queries
		qryCtrl, isQueried := ctrl.(xontrols.Queries)

		if isQueried {
			for qkey, qval := range qryCtrl.AcceptsQuery() {
				ctrlSub.Queries(qkey, qval)
			}
		}

		menuPath := ctrlPath

		//more than just a slash
		if strings.HasPrefix(path, "/") && len(path) > 1 {
			menuPath = path + ctrlPath
		}

		menu = append(menu, bodies.NewItem("", menuPath, ctrlName, isCreate, nil))
	}

	e.sideMenu.AddGroup(strings.Title(strings.Replace(path, "/", "", -1)), menu)
}

func (e *colourEpoxy) filter(name string, required roletype.Enum, mxFunc mix.InitFunc, process ServeFunc) http.HandlerFunc {
	srv := e.service
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, srv.ID)

		allow, avoc := filters.TokenCookieCheck(ctx, required, srv.PublicKey, srv.Name)
		if !allow {
			err := sendToLogin(ctx, e.securityUrl)

			if err != nil {
				log.Panicln(err)
			}

			return
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := process(ctx)
		mxer := mxFunc(name, data, e.identity, avoc).(mix.ColourMixer)
		//mxer.ApplySettings(route.Name, *e.settings, avoc)

		mxer.CreateSideMenu(e.sideMenu)
		err := ctx.Serve(status, mxer)

		if err != nil {
			log.Panicln(err)
		}
	}
}

func getControllerName(ctrl xontrols.Nomad) string {
	tpe := reflect.TypeOf(ctrl).String()
	lstDot := strings.LastIndex(tpe, ".")

	if lstDot != -1 {
		return tpe[(lstDot + 1):]
	}

	return tpe
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
