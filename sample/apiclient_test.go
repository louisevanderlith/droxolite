package sample

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/louisevanderlith/droxolite/mix"

	"github.com/louisevanderlith/droxolite/element"

	"github.com/louisevanderlith/droxolite/resins"

	"github.com/louisevanderlith/droxolite/sample/clients"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/security/impl"
	"github.com/louisevanderlith/droxolite/security/models"
	"github.com/louisevanderlith/droxolite/security/roletype"
)

func init() {
	host := ".localhost/"
	uri := ".localhost:8090/"
	clnt := models.NewPrivateClient("TestAPI", "Sample Api Client", uri)
	clnt.Scopes = append(clnt.Scopes, models.NewScope("testapi", "Test API", "just a scope to test with", []models.Claim{
		{"api:user", "Allows User base interactions with API"},
	}))

	cs := impl.NewFakeRegister()
	cred, err := cs.AddClient(clnt)

	if err != nil {
		panic(err)
	}

	intro := impl.NewFakeInspector(cs)
	apiEpoxy = resins.NewMonoEpoxy(cred, intro, element.GetNoTheme(host, cred.ID, ""))
	apiRoutes(apiEpoxy)
	apiEpoxy.EnableCORS(host)
}

func TestNomad_GetAcceptsQuery_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/nomad?name=Jannie", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	result := ""
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err, rr.Body.String())
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := "Nomad got Jannie"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestStore_Get_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/store", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	t.Log(rr.Body.String())

	var result []string
	rest, err := bodies.MarshalToResult(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	if len(rest.Reason) > 0 {
		t.Fatalf(rest.Reason)
	}

	expected := []string{"Berry", "Orange", "Apple"}

	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Fatalf("unexpected body: got %v want %v", result[i], expected[i])
		}
	}
}

func TestStore_GetOne_OK(t *testing.T) {
	rr, err := GetResponse(apiEpoxy, "/fake/store/1560674025%601", nil)

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

	expected := "Got a Key 1560674025`1"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestStore_Create_OK(t *testing.T) {

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
	rr, err := GetResponse(apiEpoxy, "/fake/store/C78/eyJuYW1lIjogIkppbW15IiwiYWdlOiB7ICJtb250aCI6IDIsICJkYXRlIjogOCwgInllYXIiOiAxOTkxfSwiYWxpdmUiOiB0cnVlfQ==", nil)

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
	rr, err := GetResponse(apiEpoxy, "/fake/store/73", readr)

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

	expected := "#73: Jump"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func apiRoutes(e resins.Epoxi) {
	e.JoinBundle("/fake", roletype.Unknown, mix.JSON, &clients.Nomad{}, &clients.Store{}, &clients.Apx{})

	/*fakeCtrl := &FakeAPI{}

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
	poxy.JoinBundle("Fake", roletype.Unknown, fakeCtrl)
	poxy.AddBundle(fkgroup)

	subCtrl := &sub.SubAPICtrl{}
	subGroup := routing.NewRouteGroup("Sub", mix.JSON)
	subGroup.AddRoute("Sub Home", "", http.MethodGet, roletype.Unknown, subCtrl.Get)

	complxCtrl := &sub.ComplexAPICtrl{}
	complxGroup := routing.NewRouteGroup("Complex", mix.JSON)
	complxGroup.AddRoute("Sub Complex Home", "", http.MethodGet, roletype.Unknown, complxCtrl.Get)

	subGroup.AddSubGroup(complxGroup)
	poxy.AddBundle(subGroup)*/
}
