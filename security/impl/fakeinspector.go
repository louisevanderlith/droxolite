package impl

import (
	"encoding/json"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/louisevanderlith/droxolite/security/client"
	"github.com/louisevanderlith/droxolite/security/models"
)

type fakeInspector struct {
	clientStore client.Registar
}

func NewFakeInspector(clientStore client.Registar) client.Inspector {
	return &fakeInspector{clientStore}
}

func (i *fakeInspector) Introspect(token string, creds models.ClientCred) (*models.ClaimIdentity, error) {
	//err := authenticateFakeClient(client)

	_, err := i.clientStore.Authenticate(creds)

	if err != nil {
		return nil, err
	}

	return getFakeIdentity(token)
}

func getFakeIdentity(token string) (*models.ClaimIdentity, error) {
	if len(token) == 0 {
		return nil, errors.New("token empty")
	}

	claims := make(jwt.MapClaims)
	claims["name"] = "Fake User"
	claims["email"] = "fake@mango.avo"
	claims["role"] = "folio:user"
	claims["role"] = "testapi:user"
	claims["role"] = "testapp:user"

	// TODO: Assign approved claims

	iden := &jwt.Token{
		Valid:  true,
		Claims: claims,
	}

	if !iden.Valid {
		return nil, errors.New("token invalid")
	}

	jClaim, err := json.Marshal(iden.Claims)

	if err != nil {
		return nil, err
	}

	result := &models.ClaimIdentity{}
	err = json.Unmarshal(jClaim, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
