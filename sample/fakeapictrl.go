package sample

import (
	"fmt"
	"strconv"

	"github.com/louisevanderlith/husk"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type FakeAPICtrl struct {
	xontrols.APICtrl
}

// /
func (c *FakeAPICtrl) Get() {
	c.Ctx.WriteResponse([]byte("Fake GET Working"))
}

func (c *FakeAPICtrl) GetKey() {
	param := c.Ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		c.Ctx.WriteResponse([]byte(err.Error()))
		return
	}
	c.Ctx.WriteResponse([]byte(fmt.Sprintf("Got a Key %s", result)))
}

//:id
func (c *FakeAPICtrl) GetId() {
	param := c.Ctx.FindParam("id")
	result := fmt.Sprintf("We Found %v", param)
	c.Ctx.WriteResponse([]byte(result))
}

//name:/id:
func (c *FakeAPICtrl) GetName() {
	param := c.Ctx.FindParam("id")
	name := c.Ctx.FindParam("name")
	result := fmt.Sprintf("%s is %v", name, param)
	c.Ctx.WriteResponse([]byte(result))
}

//yes:
func (c *FakeAPICtrl) GetAnswer() {
	param := c.Ctx.FindParam("yes")
	yes, err := strconv.ParseBool(param)

	if err != nil {
		c.Ctx.WriteResponse([]byte(err.Error()))
		return
	}

	if !yes {
		c.Ctx.WriteResponse([]byte("Thanks for Nothing!"))
		return
	}

	c.Ctx.WriteResponse([]byte("That's Nice"))
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
