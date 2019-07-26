package droxolite_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/louisevanderlith/droxolite/context"
)

func TestSetHeader_ResponseMustHaveValue(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	ctx := context.New(resp, req)

	ctx.SetHeader("Must", "Persist")

	val, ok := resp.HeaderMap["Must"]

	if !ok {
		t.Fatal("Header not found")
	}

	if len(val) == 0 {
		t.Fatal("No values found")
	}

	if val[0] != "Persist" {
		t.Errorf("Expected 'Persist', got %s", val)
	}
}

func TestFindParam_GetValueParamFromRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/?find=me", nil)

	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	ctx := context.New(resp, req)

	val := ctx.FindParam("find")

	if len(val) == 0 {
		t.Fatal("Param not found")
	}

	if val != "me" {
		t.Errorf("Expected 'me', got %s", val)
	}
}
