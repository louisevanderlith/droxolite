package sample

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/droxolite/sample/clients"

	"github.com/louisevanderlith/droxolite/element"
	"github.com/louisevanderlith/droxolite/resins"
	"github.com/louisevanderlith/droxolite/security/impl"
	"github.com/louisevanderlith/droxolite/security/models"
	"github.com/louisevanderlith/droxolite/security/roletype"
)

var (
	appEpoxy resins.Epoxi
)

func init() {
	host := ".localhost/"
	uri := ".localhost:8091/"
	callb := "https://localhost:48091/"
	clnt := models.NewPublicClient("Test.APP", "Sample App Client", uri, callb)
	clnt.LogoURI = uri + "logo.jpg"
	clnt.Scopes = append(clnt.Scopes, models.NewScope("testapp", "Test APP", "just a scope to test with", []models.Claim{
		{"app:user", "Allows User base interactions with APP"},
	}))

	cs := impl.NewFakeRegister()
	cred, err := cs.AddClient(clnt)

	if err != nil {
		panic(err)
	}

	intro := impl.NewFakeInspector(cs)

	theme := element.GetNoTheme(host, cred.ID, "none")

	err = theme.LoadTemplate("./views", "master.html")

	if err != nil {
		panic(err)
	}

	appEpoxy = resins.NewColourEpoxy(cred, intro, theme, "auth.localhost", roletype.Unknown, clients.Index)
	appRoutes(appEpoxy)
}

func TestAPP_DistAsset_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/dist/site.css", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appEpoxy.Router()

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

	handle := appEpoxy.Router()

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

	handle := appEpoxy.Router()

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

	handle := appEpoxy.Router()

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

func appRoutes(e resins.Epoxi) {
	e.JoinBundle("/", roletype.Unknown, mix.Page, &clients.Interface{})
	e.JoinBundle("/stock", roletype.Unknown, mix.Page, &clients.Parts{}, &clients.Services{})
	/*fakeCtrl := &FakeAPP{}
	grp := routing.NewInterfaceBundle("", roletype.Unknown, fakeCtrl)
	grp.RouteGroup().AddRoute("Home", "/broken", "GET", roletype.Unknown, fakeCtrl.GetBroken)

	poxy.AddBundle(grp)

	stockGrp := routing.NewInterfaceBundle("Stock", roletype.Unknown, &Parts{}, &Services{})
	poxy.AddBundle(stockGrp)*/
}
