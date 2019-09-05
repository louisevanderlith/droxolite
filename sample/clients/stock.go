package clients

import (
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type Parts struct {
}

func (req *Parts) Get(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (req *Parts) Search(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (req *Parts) View(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (req *Parts) Create(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

type Services struct {
}

func (req *Services) Get(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (req *Services) Search(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (req *Services) View(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (req *Services) Create(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}
