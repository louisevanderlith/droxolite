package sample

import (
	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/louisevanderlith/droxolite/sample/clients"
)

func TestAPP_DistAsset_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/dist/site.css", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appRoutes()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Not OK: %v", rr.Code)
	}

	expected := "h1{margin: auto;}"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAPP_Home_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appRoutes()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Not OK: %v", rr.Code)
	}

	expected := "<h1>MasterPage</h1><p>This is the Home Page</p><p>Welcome</p>  <span>Footer</span>"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAPP_SubDefault_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/stock/parts", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appRoutes()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Not OK: %v", rr.Code)
	}

	expected := "<h1>MasterPage</h1><h1>Parts</h1>  <span>Footer</span>"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAPP_Error_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/broken", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appRoutes()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Not OK: %v", rr.Code)
	}

	expected := "<h1>MasterPage</h1><h1>Something unexptected Happended:</h1><p>this path must break</p>  <span>Footer</span>"

	if rr.Body.Len() != len(expected) {
		t.Errorf("unexpected length: got %v want %v",
			rr.Body.Len(), len(expected))
	}

	if strings.Compare(rr.Body.String(), expected) != 0 {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func appRoutes() http.Handler {
	mstr, tmpl, err := droxolite.LoadTemplate("./views", "master.html")

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	r.HandleFunc("/", clients.InterfaceGet(mstr, tmpl)).Methods(http.MethodGet)
	r.HandleFunc("/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", clients.InterfaceSearch(mstr, tmpl)).Methods(http.MethodGet)
	r.HandleFunc("/{pagesize:[A-Z][0-9]+}", clients.InterfaceSearch(mstr, tmpl)).Methods(http.MethodGet)
	r.HandleFunc("/{key:[0-9]+\x60[0-9]+}", clients.InterfaceView(mstr, tmpl)).Methods(http.MethodGet)
	r.HandleFunc("/create", clients.InterfaceCreate(mstr, tmpl)).Methods(http.MethodPost)

	stck := r.PathPrefix("/stock").Subrouter()
	stck.HandleFunc("/parts", clients.PartsGet(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", clients.PartsSearch(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/{pagesize:[A-Z][0-9]+}", clients.PartsSearch(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/{key:[0-9]+\x60[0-9]+}", clients.PartsView(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/create", clients.PartsCreate(mstr, tmpl)).Methods(http.MethodPost)
	stck.HandleFunc("/services", clients.ServicesGet(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/services/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", clients.ServicesSearch(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/services/{pagesize:[A-Z][0-9]+}", clients.ServicesSearch(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/services/{key:[0-9]+\x60[0-9]+}", clients.ServicesView(mstr, tmpl)).Methods(http.MethodGet)
	stck.HandleFunc("/services/create", clients.ServicesCreate(mstr, tmpl)).Methods(http.MethodPost)

	return r
}
