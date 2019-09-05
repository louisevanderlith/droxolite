package sample

import (
	"net/http"
	"testing"
)

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
