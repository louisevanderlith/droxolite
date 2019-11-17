package impl

import (
	"errors"
	"fmt"

	"github.com/louisevanderlith/droxolite/security/client"

	"github.com/louisevanderlith/droxolite/security/models"
	uuid "github.com/nu7hatch/gouuid"
)

type fakeRegister struct {
	clients map[string]models.Client
}

func NewFakeRegister() client.Registar {
	result := &fakeRegister{
		clients: make(map[string]models.Client),
	}

	c := models.NewPrivateClient("stock", "stock api private client", "https://stock.mango.avo")
	c.Secret = "JFmXTWvaHO"

	result.clients["tJVgDxKa"] = c

	return result
}

func (r *fakeRegister) AddClient(c models.Client) (models.ClientCred, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return models.ClientCred{}, err
	}

	u4, err := uuid.NewV4()

	if err != nil {
		return models.ClientCred{}, err
	}

	c.Secret = fmt.Sprintf("%x", u4)

	valid, err := c.Valid()

	if err != nil {
		return models.ClientCred{}, err
	}

	if !valid {
		return models.ClientCred{}, errors.New("client is invalid")
	}

	r.clients[id.String()] = c

	return models.ClientCred{id.String(), c.Secret}, nil
}

func (cs *fakeRegister) Find(ID string) (models.Client, error) {
	clnt := cs.clients[ID]

	return clnt, nil
}

func (cs *fakeRegister) FindByName(name string) (models.Client, error) {
	for _, v := range cs.clients {
		if v.Name == name {
			return v, nil
		}
	}

	return models.Client{}, errors.New("client not authorised")
}

func (cs *fakeRegister) Authenticate(cred models.ClientCred) (models.Client, error) {
	clnt, err := cs.Find(cred.ID)

	if err != nil {
		return clnt, err
	}

	if clnt.Secret != cred.Secret {
		return models.Client{}, errors.New("client not authorised")
	}

	return clnt, nil
}
