package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/louisevanderlith/droxolite/security/models"
)

type TrueRegister struct {
	routerURI string
}

func NewTrueRegister(routerURI string) *TrueRegister {
	return &TrueRegister{routerURI}
}

func (r *TrueRegister) AddClient(c models.Client) (models.ClientCred, error) {
	bits, err := json.Marshal(c)

	if err != nil {
		return models.ClientCred{}, err
	}

	disco := fmt.Sprintf("%sdiscovery", r.routerURI)
	resp, err := http.Post(disco, "application/json", bytes.NewBuffer(bits))

	if err != nil {
		return models.ClientCred{}, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return models.ClientCred{}, err
	}

	cred := models.ClientCred{}
	err = json.Unmarshal(contents, &cred)

	if err != nil {
		return models.ClientCred{}, err
	}

	return cred, nil
}
