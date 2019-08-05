package sample

import (
	"errors"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type FakeAPPCtrl struct {
	xontrols.UICtrl
}

func (c *FakeAPPCtrl) Default() {
	c.Setup("home", "Fake Home", false)
	data := "Welcome"
	err := c.Serve(http.StatusOK, nil, data)

	if err != nil {
		log.Println(err)
	}
}

func (c *FakeAPPCtrl) GetBroken() {
	c.Setup("home", "Fake Home", false)

	c.Serve(http.StatusInternalServerError, errors.New("this path must break"), nil)
}
