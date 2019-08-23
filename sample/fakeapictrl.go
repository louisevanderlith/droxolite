package sample

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/louisevanderlith/husk"

	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/mix"
)

type FakeAPI struct {
}

func (ctrl *FakeAPI) Get(ctx context.Contexer) (int, interface{}) {
	return http.StatusOK, "Fake GET Working"
}

func (c *FakeAPI) GetHash(ctx context.Contexer) (int, interface{}) {
	hsh := ctx.FindParam("hash")

	decoded, err := base64.StdEncoding.DecodeString(hsh)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, string(decoded)
}

func (c *FakeAPI) GetKey(ctx context.Contexer) (int, interface{}) {
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, fmt.Sprintf("Got a Key %s", result)
}

func (c *FakeAPI) GetPage(ctx context.Contexer) (int, interface{}) {
	page, size := ctx.GetPageData()

	return http.StatusOK, fmt.Sprintf("Page %v, Size %v", page, size)
}

//:id
func (c *FakeAPI) GetId(ctx context.Contexer) (int, interface{}) {
	param := ctx.FindParam("id")
	result := fmt.Sprintf("We Found %v", param)

	return http.StatusOK, mix.JSON(result)
}

func (c *FakeAPI) GetQueryStr(ctx context.Contexer) (int, interface{}) {
	param := ctx.FindQueryParam("name")
	result := fmt.Sprintf("Fake Query %s", param)
	return http.StatusOK, result
}

//name:/id:
func (c *FakeAPI) GetName(ctx context.Contexer) (int, interface{}) {
	param := ctx.FindParam("id")
	name := ctx.FindParam("name")
	result := fmt.Sprintf("%s is %v", name, param)
	return http.StatusOK, result
}

//yes:
func (c *FakeAPI) GetAnswer(ctx context.Contexer) (int, interface{}) {
	param := ctx.FindParam("yes")
	yes, err := strconv.ParseBool(param)

	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	if !yes {
		return http.StatusOK, "Thanks for Nothing!"
	}

	return http.StatusOK, "That's Nice"
}

// :id {string}
func (c *FakeAPI) Post(ctx context.Contexer) (int, interface{}) {
	param := ctx.FindParam("id")
	body := struct{ Act string }{}
	err := ctx.Body(&body)

	if err != nil {
		return http.StatusBadRequest, err
	}

	result := fmt.Sprintf("#%v: %s", param, body.Act)

	return http.StatusOK, result
}
