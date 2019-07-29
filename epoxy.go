package droxolite

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/xontrols"
)

const (
	readTimeout  = time.Second * 15
	writeTimeout = time.Second * 15
)

type Route struct {
	Path         string
	Method       string
	RequiredRole roletype.Enum
	Function     func()
}

type RouteGroup struct {
	Name       string
	Controller xontrols.Controller
	Routes     []*Route
}

func NewRouteGroup(name string, ctrl xontrols.Controller) *RouteGroup {
	return &RouteGroup{
		Name:       name,
		Controller: ctrl,
	}
}

func (g *RouteGroup) AddRoute(path, method string, requiredRole roletype.Enum, function func()) *Route {
	result := &Route{
		Path:         path,
		Method:       method,
		RequiredRole: requiredRole,
		Function:     function,
	}

	g.Routes = append(g.Routes, result)

	return result
}

//Epoxy puts everything together
type Epoxy struct {
	service    *Service
	router     *mux.Router
	server     *http.Server
	settings   *bodies.ThemeSetting
	sideMenu   *bodies.Menu
	masterpage string
}

//NewExpoxy returns a new Instance of the Epoxy
func NewEpoxy(service *Service) *Epoxy {
	return &Epoxy{
		service:  service,
		router:   mux.NewRouter(),
		settings: nil,
	}
}

//NewExpoxy returns a new Instance of the Epoxy with a Theme
func NewColourEpoxy(service *Service, settings bodies.ThemeSetting, masterpage string) *Epoxy {
	return &Epoxy{
		service:    service,
		router:     mux.NewRouter(),
		settings:   &settings,
		sideMenu:   bodies.NewMenu(),
		masterpage: masterpage,
	}
}

func (e *Epoxy) AddGroup(routeGroup *RouteGroup) {
	uiCtrl, isUI := routeGroup.Controller.(xontrols.UIController)

	if isUI {
		if e.settings == nil {
			log.Fatalf("Use the Colour Epoxy!")
		}

		uiCtrl.SetTheme(*e.settings, e.masterpage)

		children := bodies.NewMenu()
		for _, v := range routeGroup.Routes {
			if v.Method == "GET" {
				children.AddItem(v.Path, reflect.TypeOf(v.Function).Name(), "fa-ban", nil)
			}
		}

		e.sideMenu.AddItem("#", routeGroup.Name, "fa-home", children)
	}

	sub := e.router.PathPrefix("/" + strings.ToLower(routeGroup.Name)).Subrouter()

	for _, v := range routeGroup.Routes {
		sub.Handle(v.Path, e.Handle(routeGroup.Controller, v.Function)).Methods(v.Method)
	}
}

/*
r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        t, err := route.GetPathTemplate()
        if err != nil {
            return err
        }
        // p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
        // for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
        p, err := route.GetPathRegexp()
        if err != nil {
            return err
        }
        m, err := route.GetMethods()
        if err != nil {
            return err
        }
        fmt.Println(strings.Join(m, ","), t, p)
        return nil
    })
*/

func (e *Epoxy) Handle(ctrl xontrols.Controller, call func()) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctrl.CreateInstance(context.New(resp, req), e.service.ID)
		ctrl.Prepare()

		uiCtrl, isUI := ctrl.(xontrols.UIController)

		if isUI {
			uiCtrl.CreateSideMenu(e.sideMenu)
		}

		if ctrl.Filter() {
			call()
		} else {
			ctrl.Serve(http.StatusUnauthorized, errors.New("Not allowed"), nil)
		}
	}
}

func (e *Epoxy) GetRouter() *mux.Router {
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

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	e.server = newServer(e.service.Port)
	e.server.TLSConfig = cfg
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

/*
func serveHTTP2(router *mux.Router, httpsPort int, certPath, publicKey, privateKey string) {
	publicKeyPem := readBlocks(path.Join(certPath, publicKey))
	privateKeyPem := readBlocks(path.Join(certPath, privateKey))
	cert, err := tls.X509KeyPair(publicKeyPem, privateKeyPem)

	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	srv := &http.Server{
		TLSConfig:    cfg,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		Addr:         fmt.Sprintf(":%v", httpsPort),
		Handler:      router,
	}

	log.Println("Listening...")

	err = srv.ListenAndServeTLS("", "")

	if err != nil {
		panic(err)
	}
}

*/

/*
hosts := routers.SetupRouter(instanceID, certPath)
	//subs := domains.RegisterSubdomains(instanceID, certPath)

	go serveHTTP2(hosts, httpsPort, certPath, publicKey, privateKey)

	err := http.ListenAndServe(fmt.Sprintf(":%v", httpPort), http.HandlerFunc(redirectTLS))

*/

func (e *Epoxy) Plak() {
	//avoc, err := bodies.GetAvoCookie(ctrl.GetMyToken(), ctrl.ctrlMap.GetPublicKeyPath())

	//Add Ctx
}

/*
//Add is used to specify the permissions required for a controller's actions.
func (m *Epoxy) Add(path string, actionMap map[string]int) {
	m.mapping[path] = actionMap
}

//GetRequiredRole will return the RoleType required to access the 'path' and 'action'
func (m *Epoxy) GetRequiredRole(path, action string) (roletype.Enum, error) {
	actionMap, hasCtrl := m.mapping[path]

	if !hasCtrl {
		for actPath, actMap := range m.mapping {
			if strings.Contains(path, actPath) {
				actionMap = actMap
				break
			}
		}
	}

	if actionMap == nil {
		return roletype.Unknown, fmt.Errorf("missing mapping for %s on %s", action, path)
	}

	roleType, hasAction := actionMap[strings.ToUpper(action)]

	if !hasAction {
		return roletype.Unknown, nil
	}

	return roleType, nil
}
*/
func readBlocks(filePath string) []byte {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
