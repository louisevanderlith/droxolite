package sample

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/louisevanderlith/droxolite/element"
	"github.com/louisevanderlith/droxolite/mix"

	"github.com/louisevanderlith/droxolite/resins"
	"github.com/louisevanderlith/droxolite/routing"

	"github.com/louisevanderlith/droxolite/sample/sub"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/droxolite/servicetype"
)

var (
	apiEpoxy resins.Epoxi
)

func init() {
	srvc := bodies.NewService("Test.API", "/certs/none.pem", 8090, servicetype.API)
	srvc.ID = "Tester1"

	apiEpoxy = resins.NewBasicEpoxy(srvc, element.GetNoTheme(".localhost/", srvc.ID, ""))
	apiRoutes(apiEpoxy)
	apiEpoxy.EnableCORS(".localhost/")
}

func TestPrepare_MustHaveHeader_StrictTransportSecurity(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	val := rr.Header().Get("Strict-Transport-Security")

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}

func TestPrepare_MustHaveHeader_AccessControlAllowCredentialls(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	val := rr.Header().Get("Access-Control-Allow-Credentials")

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}

func TestPrepare_MustHaveHeader_Server(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	val := rr.Header().Get("Server")

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}

func TestPrepare_MustHaveHeader_XContentTypeOptions(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	val := rr.Header().Get("X-Content-Type-Options")

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}

/*
func TestAPI_OPTIONS_CORS(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", "/fake", nil)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Access-Control-Request-Method", "POST")           // needs to be non-empty
	req.Header.Set("Access-Control-Request-Headers", "Authorization") // needs to be non-empty
	req.Header.Set("Origin", "https://tester.localhost/")             // needs to be non-empty

	handle := apiEpoxy.Router()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	t.Log(rr.Header())

	if len(rr.Header().Get("Access-Control-Allow-Method")) == 0 {
		t.Fatal("Allow Methods not Found")
	}

	if len(rr.Header().Get("Access-Control-Allow-Origin")) == 0 {
		t.Fatal("Allow Origin not Found")
	}
}*/

func TestMain_API_DefaultPath_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Fake GET Working"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_SubPath_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/sub", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "I am a sub controller"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_SubComplexPath_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/sub/complex", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	log.Println(rr.Body.String())
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "This is complex!"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_QueryPath_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/query?name=%60", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Fake Query `"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_IdParam_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/73", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	t.Log(rr.Body.String())

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "We Found 73"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_NameAndIdParam_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/Jimmy/73", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Jimmy is 73"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v", result, expected)
	}
}

func TestMain_API_HuskKey_Escaped_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/1560674025%601", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Got a Key 1560674025`1"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_HuskKey_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/1563985947336`12", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Got a Key 1563985947336`12"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_PageSize_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/all/C78", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Page 3, Size 78"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_BooleanParam_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/question/false", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Thanks for Nothing!"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

/*

 */

func TestMain_API_HashParam_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/base/eyJuYW1lIjogIkppbW15IiwiYWdlOiB7ICJtb250aCI6IDIsICJkYXRlIjogOCwgInllYXIiOiAxOTkxfSwiYWxpdmUiOiB0cnVlfQ==", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := `{"name": "Jimmy","age: { "month": 2, "date": 8, "year": 1991},"alive": true}`
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
	rr, err := GetResponse(apiEpoxy, "/fake/73", readr)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "#73: Jump"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func apiRoutes(poxy resins.Epoxi) {
	fakeCtrl := &FakeAPI{}

	fkgroup := routing.NewRouteGroup("Fake", mix.JSON)
	fkgroup.AddRoute("Home", "", "GET", roletype.Unknown, fakeCtrl.Get)

	q := make(map[string]string)
	q["name"] = "{name}"
	fkgroup.AddRouteWithQueries("Query String", "/query", "GET", roletype.Unknown, q, fakeCtrl.GetQueryStr)
	fkgroup.AddRoute("Key", "/{key:[0-9]+\x60[0-9]+}", "GET", roletype.Unknown, fakeCtrl.GetKey)
	fkgroup.AddRoute("Id POST", "/{id:[0-9]+}", "POST", roletype.Unknown, fakeCtrl.Post)
	fkgroup.AddRoute("Id", "/{id:[0-9]+}", "GET", roletype.Unknown, fakeCtrl.GetId)
	fkgroup.AddRoute("Question Answer", "/question/{yes:true|false}", "GET", roletype.Unknown, fakeCtrl.GetAnswer)
	fkgroup.AddRoute("Name", "/{name:[a-zA-Z]+}/{id:[0-9]+}", "GET", roletype.Unknown, fakeCtrl.GetName)
	fkgroup.AddRoute("Page", "/all/{pagesize:[A-Z][0-9]+}", "GET", roletype.Unknown, fakeCtrl.GetPage)
	fkgroup.AddRoute("base", "/base/{hash:[a-zA-Z0-9]+={0,2}}", "GET", roletype.Unknown, fakeCtrl.GetHash)
	poxy.AddBundle(fkgroup)

	subCtrl := &sub.SubAPICtrl{}
	subGroup := routing.NewRouteGroup("Sub", mix.JSON)
	subGroup.AddRoute("Sub Home", "", http.MethodGet, roletype.Unknown, subCtrl.Get)

	complxCtrl := &sub.ComplexAPICtrl{}
	complxGroup := routing.NewRouteGroup("Complex", mix.JSON)
	complxGroup.AddRoute("Sub Complex Home", "", http.MethodGet, roletype.Unknown, complxCtrl.Get)

	subGroup.AddSubGroup(complxGroup)
	poxy.AddBundle(subGroup)
}
