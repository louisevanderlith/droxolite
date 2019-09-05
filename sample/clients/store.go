package clients

import (
	"fmt"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/husk"
)

type Store struct {
}

func (x *Store) Get(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, []string{"Berry", "Orange", "Apple"}
}

func (x *Store) GetOne(ctx context.Requester) (int, interface{}) {
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, fmt.Sprintf("Got a Key %s", result)
}

func (x *Store) Create(ctx context.Requester) (int, interface{}) {
	param := ctx.FindParam("id")
	body := struct{ Act string }{}
	err := ctx.Body(&body)

	if err != nil {
		return http.StatusBadRequest, err
	}

	result := fmt.Sprintf("#%v: %s", param, body.Act)

	return http.StatusOK, result
}

func (x *Store) Update(ctx context.Requester) (int, interface{}) {
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, fmt.Sprintf("Updated item with Key %s", result)
}

func (x *Store) Delete(ctx context.Requester) (int, interface{}) {
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, fmt.Sprintf("Deleted item with Key %s", result)
}
