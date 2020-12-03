package sample

import (
	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/drx"
	"github.com/louisevanderlith/droxolite/mix"
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

	expected := "<h1>MasterPage</h1><p>This is the Home Page</p><p>Welcome</p><span>Footer</span>"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func BenchmarkAPP_SubDefault_OK(b *testing.B) {
	req, err := http.NewRequest("GET", "/stock/parts", nil)

	if err != nil {
		b.Fatal(err)
	}

	handle := appRoutes()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		b.Errorf("Not OK: %v", rr.Code)
	}

	expected := "<h1>MasterPage</h1><h2>Layout X</h2><div><p>This is the parts page</p></div><span>Footer</span>"
	act := rr.Body.String()
	if act != expected {
		b.Errorf("unexpected body: got %v want %v", act, expected)
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

	expected := "<h1>MasterPage</h1><h2>Layout X</h2><div><p>This is the parts page</p></div><span>Footer</span>"
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

	expected := "<h1>MasterPage</h1><h1>Something unexpected Happened:</h1><p>this path must break</p>  <span>Footer</span>"

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
	tmpl, err := drx.LoadTemplate("./views")

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	distPath := http.FileSystem(http.Dir("dist/"))
	fs := http.FileServer(distPath)
	r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", fs))

	fact := mix.NewPageFactory(tmpl)
	r.HandleFunc("/", clients.InterfaceGet(fact)).Methods(http.MethodGet)
	r.HandleFunc("/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", clients.InterfaceSearch(fact)).Methods(http.MethodGet)
	r.HandleFunc("/{pagesize:[A-Z][0-9]+}", clients.InterfaceSearch(fact)).Methods(http.MethodGet)
	r.HandleFunc("/{key:[0-9]+\x60[0-9]+}", clients.InterfaceView(fact)).Methods(http.MethodGet)
	r.HandleFunc("/create", clients.InterfaceCreate(fact)).Methods(http.MethodPost)

	stck := r.PathPrefix("/stock").Subrouter()
	stck.HandleFunc("/parts", clients.PartsGet(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", clients.PartsSearch(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/{pagesize:[A-Z][0-9]+}", clients.PartsSearch(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/{key:[0-9]+\x60[0-9]+}", clients.PartsView(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/parts/create", clients.PartsCreate(fact)).Methods(http.MethodPost)
	stck.HandleFunc("/services", clients.ServicesGet(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/services/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", clients.ServicesSearch(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/services/{pagesize:[A-Z][0-9]+}", clients.ServicesSearch(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/services/{key:[0-9]+\x60[0-9]+}", clients.ServicesView(fact)).Methods(http.MethodGet)
	stck.HandleFunc("/services/create", clients.ServicesCreate(fact)).Methods(http.MethodPost)

	return r
}
