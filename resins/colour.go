package resins

import (
	"errors"
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
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/security/client"
	"github.com/louisevanderlith/droxolite/security/filters"
	"github.com/louisevanderlith/droxolite/security/models"
	"github.com/louisevanderlith/droxolite/security/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

type colourEpoxy struct {
	client      models.ClientCred
	intro       client.Inspector
	router      http.Handler
	identity    *element.Identity
	sideMenu    *bodies.Menu
	securityUrl string
}

//NewColourExpoxy returns a new Instance of the Epoxy with a Theme
func NewColourEpoxy(client models.ClientCred, intro client.Inspector, d *element.Identity, securityUrl string, indexRole roletype.Enum, indexPage ServeFunc) Epoxi {

	e := &colourEpoxy{
		client:      client,
		intro:       intro,
		identity:    d,
		securityUrl: securityUrl,
		sideMenu:    bodies.NewMenu(),
	}

	routr := mux.NewRouter()
	routr.HandleFunc("/", e.filter("index", indexRole, mix.Page, indexPage))
	//Applications have assets in the 'dist' folder
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	routr.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	e.router = routr

	return e
}

func (e *colourEpoxy) Port() int {
	return 0
}

func (e *colourEpoxy) Router() http.Handler {
	return e.router
}

func (e *colourEpoxy) Client() models.ClientCred {
	return e.client
}

func (e *colourEpoxy) EnableCORS(host string) {
	//No Need.
}

func (e *colourEpoxy) JoinPath(r *mux.Router, path, name, method string, required roletype.Enum, mxFunc mix.InitFunc, f ServeFunc) {
	r.Handle(path, e.filter(name, required, mxFunc, f)).Methods(method)
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
			e.JoinPath(ctrlSub, "/{pagesize:[A-Z][0-9]+}", ctrlName, http.MethodGet, required, mxFunc, searchCtrl.Search)
			e.JoinPath(ctrlSub, "/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", ctrlName, http.MethodGet, required, mxFunc, searchCtrl.Search)
		}

		//View
		viewCtrl, isView := ctrl.(xontrols.Viewable)

		if isView {
			e.JoinPath(ctrlSub, "/{key:[0-9]+\x60[0-9]+}", ctrlName+"View", http.MethodGet, required, mxFunc, viewCtrl.View)
		}

		//Create
		createCtrl, isCreate := ctrl.(xontrols.Createable)

		if isCreate {
			e.JoinPath(ctrlSub, "/create", ctrlName+"Create", http.MethodGet, required, mxFunc, createCtrl.Create)
		}

		//Update - Pages can't update (View can do whatever)
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
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req, e.client, e.intro)

		p := filters.Pack{
			RequestURI:   ctx.RequestURI(),
			Token:        ctx.GetMyToken(),
			RequiredRole: required,
			ClientName:   name,
			ClientCred:   e.Client(),
			Inspector:    e.intro,
		}

		avoc, err := p.IdentifyToken()

		if err != nil {
			log.Panicln(err)
		}

		//Calls the Controller Function
		//Context should be sent to function, so no controller is needed
		status, data := process(ctx)
		mxer := mxFunc(name, data, e.identity, avoc).(mix.ColourMixer)

		mxer.CreateSideMenu(e.sideMenu)
		err = ctx.Serve(status, mxer)

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

//Returns the [TOKEN] in 'Bearer [TOKEN]'
func getAuthorizationToken(ctx context.Contexer) (string, error) {
	authHead, err := ctx.GetHeader("Authorization")

	if err != nil {
		return "", err
	}

	parts := strings.Split(authHead, " ")
	tokenType := parts[0]
	if strings.Trim(tokenType, " ") != "Bearer" {
		return "", errors.New("Bearer Authentication only")
	}

	return parts[1], nil
}
