package sample

import (
	"errors"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

type FakeAPP struct {
}

func (c *FakeAPP) Default(ctx context.Contexer) (int, interface{}) {
	//c.Setup("home", "Fake Home", false)
	data := "Welcome"
	return http.StatusOK, data
}

func (c *FakeAPP) GetBroken(ctx context.Contexer) (int, interface{}) {
	//c.Setup("home", "Fake Home", false)

	return http.StatusInternalServerError, errors.New("this path must break")
}
