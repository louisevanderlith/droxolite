package clients

import (
	"errors"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type Apx struct {
}

func (x *Apx) Get(ctx context.Requester) (int, interface{}) {
	return http.StatusInternalServerError, errors.New("this path must break")
}

func (x *Apx) Create(ctx context.Requester) (int, interface{}) {
	return http.StatusNotImplemented, nil
}
