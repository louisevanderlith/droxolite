package sample

import (
	"log"

	"github.com/louisevanderlith/droxolite/controllers"
)

type FakeAPICtrl struct {
	controllers.APICtrl
}

// /
func (c *FakeAPICtrl) Get() {
	log.Println("Get Called")
}

//:id
func (c *FakeAPICtrl) GetId() {
	param := c.Ctx.FindParam("id")

	log.Println(param)
}

func (c *FakeAPICtrl) Post() {

}
