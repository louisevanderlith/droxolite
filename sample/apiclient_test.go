package sample

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/louisevanderlith/droxolite"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/servicetype"
)

var (
	apiEpoxy *droxolite.Epoxy
)

func init() {
	srvc := droxolite.NewService("Test.API", "/certs/none.pem", 8090, servicetype.API)
	srvc.ID = "Tester1"

	apiEpoxy = droxolite.NewEpoxy(srvc)
	apiRoutes(apiEpoxy)
}

func TestMain_API_DefaultPath_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/fake/", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "Fake GET Working"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_IdParam_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/fake/73", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "We Found 73"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_NameAndIdParam_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/fake/Jimmy/73", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "Jimmy is 73"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v", result, expected)
	}
}

func TestMain_API_HuskKey_Escaped_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/fake/1560674025%601", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "Got a Key 1560674025`1"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_HuskKey_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/fake/1563985947336`12", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "Got a Key 1563985947336`12"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_PageSize_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/fake/all/C78", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "Page 3, Size 78"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_BooleanParam_OK(t *testing.T) {
	req, err := http.NewRequest("GET", "/fake/question/false", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "Thanks for Nothing!"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_POST_OK(t *testing.T) {
	body, err := json.Marshal(struct{ Act string }{"Jump"})

	if err != nil {
		t.Fatal(err)
	}

	readr := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", "/fake/73", readr)

	if err != nil {
		t.Fatal(err)
	}

	handle := apiEpoxy.GetRouter()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if rest.Code != http.StatusOK {
		t.Fatalf(rest.Reason)
	}

	expected := "#73: Jump"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func apiRoutes(poxy *droxolite.Epoxy) {
	fakeCtrl := &FakeAPICtrl{}

	fkgroup := droxolite.NewRouteGroup("Fake", fakeCtrl)
	fkgroup.AddRoute("/", "GET", roletype.Admin, fakeCtrl.Get)
	fkgroup.AddRoute("/{key:[0-9]+\x60[0-9]+}", "GET", roletype.Admin, fakeCtrl.GetKey)
	fkgroup.AddRoute("/{id:[0-9]+}", "POST", roletype.Admin, fakeCtrl.Post)
	fkgroup.AddRoute("/{id:[0-9]+}", "GET", roletype.Admin, fakeCtrl.GetId)
	fkgroup.AddRoute("/question/{yes:true|false}", "GET", roletype.Admin, fakeCtrl.GetAnswer)
	fkgroup.AddRoute("/{name:[a-zA-Z]+}/{id:[0-9]+}", "GET", roletype.Admin, fakeCtrl.GetName)
	fkgroup.AddRoute("/all/{pagesize:[A-Z][0-9]+}", "GET", roletype.Admin, fakeCtrl.GetPage)
	poxy.AddGroup(fkgroup)
}
