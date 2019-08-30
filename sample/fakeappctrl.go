package sample

import (
	"errors"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type FakeAPP struct {
}

func (c *FakeAPP) Default(ctx context.Requester) (int, interface{}) {
	data := "Welcome"
	return http.StatusOK, data
}

func (c *FakeAPP) Search(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (c *FakeAPP) View(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, "Jimmy"
}

func (c *FakeAPP) Create(ctx context.Requester) (int, interface{}) {
	return http.StatusOK, nil
}

func (c *FakeAPP) GetBroken(ctx context.Requester) (int, interface{}) {
	return http.StatusInternalServerError, errors.New("this path must break")
}
