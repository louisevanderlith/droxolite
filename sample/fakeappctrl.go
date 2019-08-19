package sample

import (
	"errors"
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type FakeAPPCtrl struct {
	xontrols.UICtrl
}

func (c *FakeAPPCtrl) Default() error {
	c.Setup("home", "Fake Home", false)
	data := "Welcome"
	return c.Serve(http.StatusOK, nil, data)
}

func (c *FakeAPPCtrl) GetBroken() error {
	c.Setup("home", "Fake Home", false)

	return c.Serve(http.StatusInternalServerError, errors.New("this path must break"), nil)
}
