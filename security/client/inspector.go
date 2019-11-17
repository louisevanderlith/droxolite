package client

import "github.com/louisevanderlith/droxolite/security/models"

//Inspector provides a token instropection endpoint to read the values in a Token
type Inspector interface {
	Introspect(token string, client models.ClientCred) (*models.ClaimIdentity, error)
}
