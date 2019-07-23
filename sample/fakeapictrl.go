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

// /poes
func (c *FakeAPICtrl) GetPOEs() {
	log.Println("Get POESe")
}

//:id
func (c *FakeAPICtrl) GetId() {
	param := c.Ctx.FindParam("id")

	log.Println(param)
}

func (c *FakeAPICtrl) Post() {

}
