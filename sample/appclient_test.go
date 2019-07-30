package sample

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/louisevanderlith/droxolite"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/servicetype"
)

var (
	appEpoxy *droxolite.Epoxy
)

func init() {
	srvc := droxolite.NewService("Test.APP", "/certs/none.pem", 8091, servicetype.APP)
	srvc.ID = "Tester2"
	theme := droxolite.GetNoTheme(".localhost/", srvc.ID, "none")
	appEpoxy = droxolite.NewColourEpoxy(srvc, theme, "master.html")
	appRoutes(appEpoxy)
}

func TestAPP_DistAsset_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/dist/site.css", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appEpoxy.GetRouter()

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

	handle := appEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Not OK: %v", rr.Code)
	}

	expected := "<h1>MasterPage</h1><p>This is the Home Page</p><p>Welcome</p>"
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

	handle := appEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Not OK: %v", rr.Code)
	}

	expected := "<h1>MasterPage</h1><h1>Something unexptected Happended:</h1><p>this path must break</p>"

	if rr.Body.Len() != len(expected) {
		t.Errorf("unexpected length: got %v want %v",
			rr.Body.Len(), len(expected))
	}

	if strings.Compare(rr.Body.String(), expected) != 0 {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAPP_Menu_Paths(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Not OK: %v", rr.Code)
	}

	expected := "<h1>MasterPage</h1><p>This is the Home Page</p><p>Welcome</p><ul><li>Home</li><li>Broken</li></ul>"
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func appRoutes(poxy *droxolite.Epoxy) {
	fakeCtrl := &FakeAPPCtrl{}
	fkgroup := droxolite.NewRouteGroup("", fakeCtrl)
	fkgroup.AddRoute("/", "GET", roletype.Unknown, fakeCtrl.GetHome)
	fkgroup.AddRoute("/broken", "GET", roletype.Unknown, fakeCtrl.GetBroken)
	poxy.AddGroup(fkgroup)
}
