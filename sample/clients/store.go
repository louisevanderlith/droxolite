package clients

import (
	"encoding/base64"
	"fmt"
	"github.com/louisevanderlith/droxolite/mix"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/husk"
)

func StoreGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(w, r)
	mxr := mix.JSON([]string{"Berry", "Orange", "Apple"})
	err := ctx.Serve(http.StatusOK, mxr)

	if err != nil {
		log.Println(err)
	}
}

func StoreSearch(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(w, r)
	page, size := ctx.GetPageData()
	hsh := ctx.FindParam("hash")

	decoded, err := base64.StdEncoding.DecodeString(hsh)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mxr := mix.JSON(fmt.Sprintf("Page: %v Size: %v Decode: %s", page, size, string(decoded)))

	err = ctx.Serve(http.StatusOK, mxr)

	if err != nil {
		log.Println(err)
	}
}

func StoreView(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(w, r)
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mxr := mix.JSON(fmt.Sprintf("Got a Key %s", result))

	err = ctx.Serve(http.StatusOK, mxr)

	if err != nil {
		log.Println(err)
	}
}

func StoreCreate(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(w, r)
	param := ctx.FindParam("id")
	body := struct{ Act string }{}
	err := ctx.Body(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mxr := mix.JSON(fmt.Sprintf("#%v: %s", param, body.Act))

	err = ctx.Serve(http.StatusOK, mxr)

	if err != nil {
		log.Println(err)
	}
}

func StoreUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(w, r)
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mxr := mix.JSON(fmt.Sprintf("Updated item with Key %s", result))

	err = ctx.Serve(http.StatusOK, mxr)

	if err != nil {
		log.Println(err)
	}
}

func StoreDelete(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(w, r)
	param := ctx.FindParam("key")
	result, err := husk.ParseKey(param)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mxr := mix.JSON(fmt.Sprintf("Deleted item with Key %s", result))

	err = ctx.Serve(http.StatusOK, mxr)

	if err != nil {
		log.Println(err)
	}
}
