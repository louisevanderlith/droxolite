package sample

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testing"

	"github.com/louisevanderlith/droxolite/sample/clients"
)

func TestStore_Get_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/store", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	t.Log(rr.Body.String())

	var result []string
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"Berry", "Orange", "Apple"}

	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Fatalf("unexpected body: got %v want %v", result[i], expected[i])
		}
	}
}

func TestStore_GetOne_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/store/1560674025%601", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
		return
	}

	t.Log(rr.Body.String())

	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := "Got a Key 1560674025`1"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_NameAndIdParam_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/Jimmy/73", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
		return
	}

	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := "Jimmy is 73"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v", result, expected)
	}
}

func TestMain_API_HuskKey_Escaped_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/1560674025%601", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := "Got a Key 1560674025`1"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_HuskKey_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/store/1563985947336`12", nil)

	if err != nil {
		t.Fatal(err)
	}
	log.Println(rr.Body.String())
	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := "Got a Key 1563985947336`12"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_PageSize_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/store/C78/eyJuYW1lIjogIkppbW15IiwiYWdlOiB7ICJtb250aCI6IDIsICJkYXRlIjogOCwgInllYXIiOiAxOTkxfSwiYWxpdmUiOiB0cnVlfQ==", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := "Page 3, Size 78"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_BooleanParam_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/question/false", nil)

	if err != nil {
		t.Fatal(err)
	}

	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := "Thanks for Nothing!"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func TestMain_API_HashParam_OK(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/store/A1/eyJuYW1lIjogIkppbW15IiwiYWdlOiB7ICJtb250aCI6IDIsICJkYXRlIjogOCwgInllYXIiOiAxOTkxfSwiYWxpdmUiOiB0cnVlfQ==", nil)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
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
	rr, err := GetResponse(apiRoutes(), "/fake/store/73", readr)

	if err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Fatal(rr.Body.String())
	}

	result := ""
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	if err != nil {
		t.Fatal(err)
	}

	expected := "#73: Jump"
	if result != expected {
		t.Errorf("unexpected body: got %v want %v",
			result, expected)
	}
}

func apiRoutes() http.Handler {
	r := mux.NewRouter()
	fke := r.PathPrefix("/fake").Subrouter()

	fke.HandleFunc("/store", clients.StoreGet).Methods(http.MethodGet)
	fke.HandleFunc("/store/{key:[0-9]+\x60[0-9]+}", clients.StoreView).Methods(http.MethodGet)
	fke.HandleFunc("/store/{pagesize:[A-Z][0-9]+}/{hash:[a-zA-Z0-9]+={0,2}}", clients.StoreSearch).Methods(http.MethodGet)
	fke.HandleFunc("/store/{pagesize:[A-Z][0-9]+}", clients.StoreSearch).Methods(http.MethodGet)
	fke.HandleFunc("/store", clients.StoreCreate).Methods(http.MethodPost)
	fke.HandleFunc("/store/{key:[0-9]+`[0-9]+}", clients.StoreUpdate).Methods(http.MethodPut)
	fke.HandleFunc("/store/{key:[0-9]+`[0-9]+}", clients.StoreDelete).Methods(http.MethodDelete)

	return r
}
