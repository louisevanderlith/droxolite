package sample

import (
	"fmt"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type FakeAPICtrl struct {
	xontrols.APICtrl
}

// /
func (c *FakeAPICtrl) Get() {
	c.Ctx.WriteResponse([]byte("Fake GET Working"))
}

//:id
func (c *FakeAPICtrl) GetId() {
	param := c.Ctx.FindParam("id")
	result := fmt.Sprintf("We Found %v", param)
	c.Ctx.WriteResponse([]byte(result))
}

// :id {string}
func (c *FakeAPICtrl) Post() {
	param := c.Ctx.FindParam("id")
	body := struct{ Act string }{}
	err := c.Ctx.Body(&body)

	if err != nil {
		panic(err)
	}

	result := fmt.Sprintf("#%v: %s", param, body.Act)

	c.Ctx.WriteResponse([]byte(result))
}
