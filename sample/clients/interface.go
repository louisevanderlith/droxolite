package clients

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/husk"
)

type Interface struct {
}

func (x *Interface) Get(ctx context.Requester) (int, interface{}) {
	data := "Welcome"
	return http.StatusOK, data
}

func (x *Interface) Search(ctx context.Requester) (int, interface{}) {
	hsh := ctx.FindParam("hash")

	decoded, err := base64.StdEncoding.DecodeString(hsh)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, string(decoded)
}

func (x *Interface) View(ctx context.Requester) (int, interface{}) {
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, fmt.Sprintf("Viewing %s", result)
}

func (x *Interface) Create(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}
