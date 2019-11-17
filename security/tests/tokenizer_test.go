package tests

import (
	"net/http"
	"testing"

	"github.com/louisevanderlith/droxolite/security/signing"

	"github.com/louisevanderlith/droxolite/security/impl"
	"github.com/louisevanderlith/droxolite/security/models"
)

func init() {
	err := signing.Initialize("certs/")

	if err != nil {
		panic(err)
	}

	testPrivKey = signing.PrivateKey
}

func TestToken_GrantPassword_ROFlow(t *testing.T) {
	us := impl.NewFakeUserStore()
	reg := impl.NewFakeRegister()
	tokn := impl.Tokening{
		GrantType:    "password",
		ClientID:     testClientID,
		ClientSecret: testClientSecret,

		Audience: "unittest",
		Scope:    "stock",
		Username: testUserName,
		Password: testUserPassword,
	}
	status, resp := tokn.GrantPassword("grantpassowrd_roflow", testPrivKey, us, reg)

	if status != http.StatusOK {
		t.Fatal(resp)
	}

	toknResp := resp.(models.TokenResponse)

	if toknResp.Type != "Bearer" {
		t.Fatalf("Expected 'Bearer', got %s", toknResp.Type)
	}

	if len(toknResp.Value) == 0 {
		t.Error("token value invalid")
	}
}

func TestToken_GrantPassword_InvalidCredentials(t *testing.T) {
	us := impl.NewFakeUserStore()
	reg := impl.NewFakeRegister()
	tokn := impl.Tokening{
		GrantType:    "password",
		ClientID:     testClientID,
		ClientSecret: testClientSecret,

		Audience: "unittest",
		Scope:    "stock",
		Username: testUserName,
		Password: "wrongPassword",
	}
	status, resp := tokn.GrantPassword("grantpassowrd_roflow", testPrivKey, us, reg)

	if status != http.StatusUnauthorized {
		t.Fatal(resp)
	}
}
