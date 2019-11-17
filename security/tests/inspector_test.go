package tests

import (
	"reflect"
	"testing"

	"github.com/louisevanderlith/droxolite/security/impl"
	"github.com/louisevanderlith/droxolite/security/models"
	"github.com/louisevanderlith/droxolite/security/roletype"
)

func TestIntrospect_ValidToken_HasClaims(t *testing.T) {
	cred := models.ClientCred{
		ID:     testClientID,
		Secret: testClientSecret,
	}

	reg := impl.NewFakeRegister()
	fi := impl.NewFakeInspector(reg)
	idn, err := fi.Introspect("token", cred)

	if err != nil {
		t.Fatal(err)
	}

	var roles []models.Role
	folioRole := models.Role{
		ApplicationName: "folio",
		Description:     roletype.User,
	}
	roles = append(roles, folioRole)

	expect := models.ClaimIdentity{
		Active:   true,
		ClientID: cred.ID,
		Audience: "unittest",
		Email:    testUserName,
		Name:     "Fake User",
		Subject:  "0",
		Roles:    roles,
	}

	if !reflect.DeepEqual(idn, expect) {
		t.Errorf("Expected %+v, got %+v", expect, idn)
	}
}

func TestIntrospect_InvalidToken_NoClaims(t *testing.T) {
	/*expect := models.ClaimIdentity{
		Active: false,
		Scope:  "",
	}*/

}
