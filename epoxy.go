package droxolite

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
	"github.com/rs/cors"
)

const (
	readTimeout  = time.Second * 15
	writeTimeout = time.Second * 15
)

//Epoxy puts everything together
type Epoxy struct {
	service    *Service
	router     http.Handler //*mux.Router
	server     *http.Server
	settings   *bodies.ThemeSetting
	sideMenu   *bodies.Menu
	masterpage string
	templates  *template.Template
}

//NewExpoxy returns a new Instance of the Epoxy
func NewEpoxy(service *Service) *Epoxy {
	routr := mux.NewRouter()

	return &Epoxy{
		service:  service,
		router:   routr,
		settings: nil,
		sideMenu: bodies.NewMenu(),
	}
}

//NewExpoxy returns a new Instance of the Epoxy with a Theme
func NewColourEpoxy(service *Service, settings bodies.ThemeSetting, masterpage string) *Epoxy {
	routr := mux.NewRouter()

	//Applications have assets in the 'dist' folder
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	routr.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	e := &Epoxy{
		service:    service,
		router:     routr,
		settings:   &settings,
		masterpage: masterpage,
		sideMenu:   bodies.NewMenu(),
	}

	err := e.settings.LoadTemplate("./views", masterpage)

	if err != nil {
		panic(err)
	}

	return e
}

//EnableCORS enables host 'https://*{.localhost/}'
func (e *Epoxy) EnableCORS(host string) {
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

func (e *Epoxy) AddNamedGroup(name string, routeGroup *RouteGroup) {
	uiCtrl, isUI := routeGroup.Controller.(xontrols.UIController)

	if isUI {
		if e.settings == nil {
			log.Fatalf("Use the Colour Epoxy!")
		}

		uiCtrl.SetTheme(*e.settings, e.masterpage)

		var menuGroup []bodies.MenuItem
		for k, v := range routeGroup.Routes {
			if v.Method == http.MethodGet {
				menuGroup = append(menuGroup, bodies.NewItem(fmt.Sprintf("r%v", k), v.Path, v.Name, nil))
				//menuGroup.AddItem(v.Path, v.Name, nil)
			}
		}

		for _, sgroup := range routeGroup.SubGroups {

			var menuChildren []bodies.MenuItem

			for k, v := range sgroup.Routes {
				if v.Method == http.MethodGet {
					menuChildren = append(menuChildren, bodies.NewItem(fmt.Sprintf("c%v", k), v.Path, v.Name, nil))
					//menuGroup.AddItem(v.Path, v.Name, nil)
				}
			}

			menuGroup = append(menuChildren)
		}

		//item := .AddItem("#", routeGroup.Name, children)

		e.sideMenu.AddGroup(name, menuGroup)
	}

	sub := e.router.(*mux.Router).PathPrefix("/" + strings.ToLower(routeGroup.Name)).Subrouter()

	for _, v := range routeGroup.Routes {
		r := sub.Handle(v.Path, e.Handle(routeGroup.Controller, v.RequiredRole, v.Function)).Methods(v.Method)

		for qkey, qval := range v.Queries {
			r.Queries(qkey, qval)
		}
	}

	//add sub groups
	for _, sgroup := range routeGroup.SubGroups {
		xsub := sub.PathPrefix("/" + strings.ToLower(sgroup.Name)).Subrouter()

		for _, v := range sgroup.Routes {
			r := xsub.Handle(v.Path, e.Handle(sgroup.Controller, v.RequiredRole, v.Function)).Methods(v.Method)

			for qkey, qval := range v.Queries {
				r.Queries(qkey, qval)
			}
		}
	}
}

func (e *Epoxy) AddGroup(routeGroup *RouteGroup) {
	e.AddNamedGroup("", routeGroup)
}

func (e *Epoxy) Handle(ctrl xontrols.Controller, requiredRole roletype.Enum, call func()) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := context.New(resp, req)
		ctrl.CreateInstance(ctx, e.service.ID)
		ctrl.Prepare()

		if !ctrl.Filter(requiredRole, e.service.PublicKey, e.service.Name) {
			err := sendToLogin(ctrl.Ctx(), e.service.ID)

			if err != nil {
				ctrl.Serve(http.StatusInternalServerError, err, nil)
			}

			return
		}

		uiCtrl, isUI := ctrl.(xontrols.UIController)

		if isUI {
			uiCtrl.CreateSideMenu(e.sideMenu)
		}

		//Calls the Controller Function
		call()
	}
}

func sendToLogin(ctx context.Contexer, instanceID string) error {
	securityURL, err := GetServiceURL(instanceID, "Auth.APP", true)

	if err != nil {
		return err
	}

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

func (e *Epoxy) GetRouter() http.Handler {
	return e.router
}

//Boot starts the Epoxy Objects to serve a configured service.
func (e *Epoxy) Boot() error {
	e.server = newServer(e.service.Port)
	e.server.Handler = e.router

	return e.server.ListenAndServe()
}

//Boot starts the Epoxy Objects to securely serve a configured service
func (e *Epoxy) BootSecure(privKeyPath string, fromPort int) error {
	publicKeyPem := readBlocks(e.service.PublicKey)
	privateKeyPem := readBlocks(privKeyPath)
	cert, err := tls.X509KeyPair(publicKeyPem, privateKeyPem)

	if err != nil {
		return err
	}

	e.server = newServer(e.service.Port)
	e.server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	e.server.Handler = e.router

	err = e.server.ListenAndServeTLS("", "")

	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%v", fromPort), http.HandlerFunc(redirectTLS))
}

func (e *Epoxy) Shutdown() {
	//e.server.Shutdown(e.)
}

func newServer(port int) *http.Server {
	return &http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         fmt.Sprintf(":%v", port),
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	moveURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	http.Redirect(w, r, moveURL, http.StatusPermanentRedirect)
}

func readBlocks(filePath string) []byte {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
