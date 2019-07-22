package droxolite_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/controllers"
)

func TestPrepare_MustHaveHeader_StrictTransportSecurity(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	ctrl := new(controllers.APICtrl)
	ctx := context.New(resp, req)
	ctrl.CreateInstance(ctx)
	ctrl.Prepare()

	val, ok := resp.HeaderMap["Strict-Transport-Security"]

	if !ok {
		t.Fatal("Header not Found")
	}

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}

func TestPrepare_MustHaveHeader_AccessControlAllowCredentialls(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	ctrl := new(controllers.APICtrl)
	ctx := context.New(resp, req)
	ctrl.CreateInstance(ctx)
	ctrl.Prepare()

	val, ok := resp.HeaderMap["Access-Control-Allow-Credentials"]

	if !ok {
		t.Fatal("Header not Found")
	}

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}

func TestPrepare_MustHaveHeader_Server(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	ctrl := new(controllers.APICtrl)
	ctx := context.New(resp, req)
	ctrl.CreateInstance(ctx)
	ctrl.Prepare()

	val, ok := resp.HeaderMap["Server"]

	if !ok {
		t.Fatal("Header not Found")
	}

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}

func TestPrepare_MustHaveHeader_XContentTypeOptions(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	ctrl := new(controllers.APICtrl)
	ctx := context.New(resp, req)
	ctrl.CreateInstance(ctx)
	ctrl.Prepare()

	val, ok := resp.HeaderMap["X-Content-Type-Options"]

	if !ok {
		t.Fatal("Header not Found")
	}

	if len(val) == 0 {
		t.Fatal("No values set")
	}
}
