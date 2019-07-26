package sample

import (
	"errors"
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type FakeAPPCtrl struct {
	xontrols.UICtrl
}

func (c *FakeAPPCtrl) GetHome() {
	c.Setup("home", "Fake Home", false)
	data := "Welcome"
	c.Serve(http.StatusOK, nil, data)
}

func (c *FakeAPPCtrl) GetBroken() {
	c.Setup("home", "Fake Home", false)

	c.Serve(http.StatusInternalServerError, errors.New("this path must break"), nil)
}
