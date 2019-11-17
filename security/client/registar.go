package client

import (
	"github.com/louisevanderlith/droxolite/security/models"
)

//Registar provides a way for Client applications to register with authorization server
type Registar interface {
	AddClient(c models.Client) (models.ClientCred, error)
	Find(ID string) (models.Client, error)
	FindByName(name string) (models.Client, error)
	Authenticate(cred models.ClientCred) (models.Client, error)
}
