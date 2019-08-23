package sample

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/louisevanderlith/droxolite"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/resins"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/routing"
	"github.com/louisevanderlith/droxolite/servicetype"
)

var (
	appEpoxy resins.Epoxi
)

func init() {
	srvc := bodies.NewService("Test.APP", "/certs/none.pem", 8091, servicetype.APP)
	srvc.ID = "Tester2"
	theme := droxolite.GetNoTheme(".localhost/", srvc.ID, "none")
	appEpoxy = resins.NewColourEpoxy(srvc, theme, "master.html", "auth.localhost")
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

func TestAPP_Menu_Paths(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appEpoxy.Router()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Not OK: %v", rr.Code)
	}

	expected := `<h1>MasterPage</h1><p>This is the Home Page</p><p>Welcome</p>
	<aside>
	<p>
			Home
		</p>
		<ul>
			<li><a href="/stock">Stock</a></li>
			<li>
				<a href="/stock/parts">Parts</a>
				<ul>
					<li><a href="/stock/parts/create">Create</a></li>
				</ul>
			</li>
			<li>
			<a href="/stock/services">Services</a>
			<ul>
				<li><a href="/stock/services/create">Create</a></li>
			</ul>
		</li>
		</ul>
		<p>
			Stock.API
		</p>
		<ul>
			<li><a href="/stock">Stock</a></li>
			<li>
				<a href="/stock/parts">Parts</a>
				<ul>
					<li><a href="/stock/parts/create">Create</a></li>
				</ul>
			</li>
			<li>
			<a href="/stock/services">Services</a>
			<ul>
				<li><a href="/stock/services/create">Create</a></li>
			</ul>
		</li>
		</ul>
	</aside>`
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func appRoutes(poxy resins.Epoxi) {
	fakeCtrl := &FakeAPP{}
	grp := routing.NewInterfaceBundle("", roletype.Unknown, fakeCtrl)
	grp.AddRoute("Home", "/broken", "GET", roletype.Unknown, fakeCtrl.GetBroken)

	poxy.AddGroup(grp)

	stockGrp := routing.NewInterfaceBundle("Stock", roletype.Unknown, &Parts{}, &Services{})
	poxy.AddGroup(stockGrp)
}
