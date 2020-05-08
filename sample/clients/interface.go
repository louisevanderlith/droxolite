package clients

import (
	"encoding/base64"
	"fmt"
	"github.com/louisevanderlith/droxolite/mix"
	"html/template"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/husk"
)

func InterfaceGet(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)

		mxr := mix.Page("index", "You're Home!", ctx.GetTokenInfo(), mstr, tmpl)

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceSearch(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		hsh := ctx.FindParam("hash")

		decoded, err := base64.StdEncoding.DecodeString(hsh)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mxr := mix.Page("index", string(decoded), ctx.GetTokenInfo(), mstr, tmpl)

		err = ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceView(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		param := ctx.FindParam("key")
		result, err := husk.ParseKey(param)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mxr := mix.Page("Index", fmt.Sprintf("Viewing %s", result), ctx.GetTokenInfo(), mstr, tmpl)

		err = ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceCreate(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := mix.Page("Index", nil, ctx.GetTokenInfo(), mstr, tmpl)

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}
