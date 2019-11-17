package filters

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/louisevanderlith/droxolite/security/client"
	"github.com/louisevanderlith/droxolite/security/models"
	"github.com/louisevanderlith/droxolite/security/roletype"
)

type Pack struct {
	RequestURI   string
	Token        string
	RequiredRole roletype.Enum
	ClientName   string //used to be service name
	ClientCred   models.ClientCred
	Inspector    client.Inspector
}

func (p Pack) IdentifyToken() (*models.ClaimIdentity, error) {

	if strings.HasPrefix(p.RequestURI, "/static") || strings.HasPrefix(p.RequestURI, "/favicon") {
		return nil, nil
	}

	if len(p.Token) == 0 && p.RequiredRole == roletype.Unknown {
		return nil, nil
	}

	cred, err := p.Inspector.Introspect(p.Token, p.ClientCred)

	if err != nil {
		return nil, err
	}

	if p.RequiredRole == roletype.Unknown {
		return cred, nil
	}

	allowed, err := IsAllowed(p.ClientName, cred.Claims(), p.RequiredRole)

	if err != nil {
		return nil, err
	}

	if !allowed {
		return nil, errors.New("unauthorized")
	}

	return cred, nil
}

func IsAllowed(clientName string, claims jwt.MapClaims, required roletype.Enum) (bool, error) {
	err := claims.Valid()
	if err != nil {
		return false, err
	}

	role, err := getRole(clientName, claims)

	if err != nil {
		return false, err
	}

	return role <= required, nil
}

func getRole(appName string, usrRoles map[string]interface{}) (roletype.Enum, error) {
	role, ok := usrRoles[appName]

	if !ok {
		msg := fmt.Errorf("application permission required. %s", appName)
		return roletype.Unknown, msg
	}

	return role.(roletype.Enum), nil
}
