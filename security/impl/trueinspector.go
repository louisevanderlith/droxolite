package impl

import (
	"errors"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/do"
	"github.com/louisevanderlith/droxolite/security/models"
)

type TrueInspector struct {
	instanceID string
}

func NewTrueInspector(instanceID string) *TrueInspector {
	return &TrueInspector{
		instanceID: instanceID,
	}
}

func (i *TrueInspector) Introspect(token string, client models.ClientCred) (*models.ClaimIdentity, error) {
	claims := &models.ClaimIdentity{}
	code, err := do.SEND(http.MethodPost, "", claims, i.instanceID, "Proof.API", "introspect", client)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if code != http.StatusOK {
		log.Println("validation failed")
		return nil, errors.New("validation failed")
	}

	return claims, nil
}
