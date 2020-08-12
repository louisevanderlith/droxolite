package clients

import (
	"encoding/base64"
	"fmt"
	"github.com/louisevanderlith/droxolite/mix"
	"html/template"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/drx"
	"github.com/louisevanderlith/husk"
)

func InterfaceGet(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Index", tmpl, "./views/index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, "You are Home!"))

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceSearch(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Index", tmpl, "./views/index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		hsh := drx.FindParam(r, "hash")

		decoded, err := base64.StdEncoding.DecodeString(hsh)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = mix.Write(w, pge.Create(r, string(decoded)))

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceView(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Index", tmpl, "./views/index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		param := drx.FindParam(r, "key")
		result, err := husk.ParseKey(param)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := fmt.Sprintf("Viewing %s", result)

		err = mix.Write(w, pge.Create(r, data))

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceCreate(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Index", tmpl, "./views/index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, "Nothing Created"))

		if err != nil {
			log.Println(err)
		}
	}
}
