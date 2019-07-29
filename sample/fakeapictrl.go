package sample

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/louisevanderlith/husk"

	"github.com/louisevanderlith/droxolite/xontrols"
)

type FakeAPICtrl struct {
	xontrols.APICtrl
}

// /
func (c *FakeAPICtrl) Get() {
	c.Serve(http.StatusOK, nil, "Fake GET Working")
}

func (c *FakeAPICtrl) GetKey() {
	param := c.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		c.Serve(http.StatusInternalServerError, err, nil)
		return
	}

	c.Serve(http.StatusOK, nil, fmt.Sprintf("Got a Key %s", result))
}

func (c *FakeAPICtrl) GetPage() {
	page, size := c.GetPageData()

	c.Serve(http.StatusOK, nil, fmt.Sprintf("Page %v, Size %v", page, size))
}

//:id
func (c *FakeAPICtrl) GetId() {
	param := c.FindParam("id")
	result := fmt.Sprintf("We Found %v", param)
	c.Serve(http.StatusOK, nil, result)
}

//name:/id:
func (c *FakeAPICtrl) GetName() {
	param := c.FindParam("id")
	name := c.FindParam("name")
	result := fmt.Sprintf("%s is %v", name, param)
	c.Serve(http.StatusOK, nil, result)
}

//yes:
func (c *FakeAPICtrl) GetAnswer() {
	param := c.FindParam("yes")
	yes, err := strconv.ParseBool(param)

	if err != nil {
		log.Println(err)
		c.Serve(http.StatusInternalServerError, err, nil)
		return
	}

	if !yes {
		c.Serve(http.StatusOK, nil, "Thanks for Nothing!")
		return
	}

	c.Serve(http.StatusOK, nil, "That's Nice")
}

// :id {string}
func (c *FakeAPICtrl) Post() {
	param := c.FindParam("id")
	body := struct{ Act string }{}
	err := c.Body(&body)

	if err != nil {
		log.Fatal(err)
	}

	result := fmt.Sprintf("#%v: %s", param, body.Act)

	c.Serve(http.StatusOK, nil, result)
}
