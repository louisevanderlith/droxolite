package sample

import (
	"net/http"
	"testing"
)

func TestPrepare_MustHaveHeader_StrictTransportSecurity(t *testing.T) {
	rr, err := GetResponse(apiRoutes(), "/fake/store", nil)

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
	rr, err := GetResponse(apiRoutes(), "/fake/store", nil)

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
	rr, err := GetResponse(apiRoutes(), "/fake/store", nil)

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
	rr, err := GetResponse(apiRoutes(), "/fake/store", nil)

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
