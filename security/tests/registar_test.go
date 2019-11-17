package tests

import (
	"testing"

	"github.com/louisevanderlith/droxolite/security/impl"

	"github.com/louisevanderlith/droxolite/security/models"
)

//Register Public Client (APP)
func TestRegisterClient_Public_CompliesWith_AuthCodeFlow(t *testing.T) {
	c := models.Client{
		Name:               "www",
		Description:        "www is an example public application client",
		RedirectURI:        "https://localhost:48091",
		EndpointAuthMethod: "none",
		GrantType:          "authorization_code",
		ResponseType:       "code",
		ClientURI:          "https://www.mango.avo",
		LogoURI:            "https://www.mango.avo/favicon.ico",
		Scopes: []models.Scope{
			models.NewScope("comment", "Comment API", "Comment allows the user to create and view comments.", []models.Claim{{"comment:user", "Allows User base interactions with Comment"}}),
		},
	}

	fr := impl.NewFakeRegister()
	creds, err := fr.AddClient(c)

	if err != nil {
		t.Fatal(err)
	}

	if len(creds.ID) == 0 {
		t.Error("expecting id")
	}

	if len(creds.Secret) == 0 {
		t.Error("expecting secret")
	}
}

//Register Private Client (API)
//Basic client authentication
func TestRegisterClient_Basic_Private_CompliesWith_ResourceOwnerPasswordCredentials(t *testing.T) {
	c := models.Client{
		Name:        "folio",
		Description: "folio is an example private webservice client",
		//No Redirect Required
		EndpointAuthMethod: "client_secret_basic",
		GrantType:          "password",
		ResponseType:       "token",
		ClientURI:          "https://folio.mango.avo",
		LogoURI:            "https://folio.mango.avo/favicon.ico",
		Scopes: []models.Scope{
			models.NewScope("folio", "Folio API", "Folio allows the user to create and view portfolia.", []models.Claim{
				{"folio:user", "Allows User base interactions with Folio"},
				{"folio:admin", "Allows Admin base interactions with Folio"},
			}),
		},
	}

	fr := impl.NewFakeRegister()
	creds, err := fr.AddClient(c)

	if err != nil {
		t.Fatal(err)
	}

	if len(creds.ID) == 0 {
		t.Error("expecting id")
	}

	if len(creds.Secret) == 0 {
		t.Error("expecting secret")
	}
}

//Register Private Client (API)
//Post cilent Authentication
func TestRegisterClient_Post_Private_CompliesWith_ResourceOwnerPasswordCredentials(t *testing.T) {
	c := models.Client{
		Name:        "folio",
		Description: "folio is an example private webservice client",
		//No Redirect Required
		EndpointAuthMethod: "client_secret_post",
		GrantType:          "password",
		ResponseType:       "token",
		ClientURI:          "https://folio.mango.avo",
		LogoURI:            "https://folio.mango.avo/favicon.ico",
		Scopes: []models.Scope{
			models.NewScope("folio", "Folio API", "Folio allows the user to create and view portfolia.", []models.Claim{
				{"folio:user", "Allows User base interactions with Folio"},
				{"folio:admin", "Allows Admin base interactions with Folio"},
			}),
		},
	}

	fr := impl.NewFakeRegister()
	creds, err := fr.AddClient(c)

	if err != nil {
		t.Fatal(err)
	}

	if len(creds.ID) == 0 {
		t.Error("expecting id")
	}

	if len(creds.Secret) == 0 {
		t.Error("expecting secret")
	}
}
