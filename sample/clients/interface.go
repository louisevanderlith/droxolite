package clients

import (
	"encoding/base64"
	"fmt"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/husk/keys"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/drx"
)

func InterfaceGet(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := mix.NewDataBag("You are Home!")
		err := mix.Write(w, fact.Create(r, "Index", "./views/index.html", b))

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceSearch(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hsh := drx.FindParam(r, "hash")

		decoded, err := base64.StdEncoding.DecodeString(hsh)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := mix.NewDataBag(string(decoded))
		err = mix.Write(w, fact.Create(r, "Index", "./views/index.html", b))

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceView(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := drx.FindParam(r, "key")
		result, err := keys.ParseKey(param)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := mix.NewDataBag(fmt.Sprintf("Viewing %s", result))

		err = mix.Write(w, fact.Create(r, "Index", "./views/index.html", b))

		if err != nil {
			log.Println(err)
		}
	}
}

func InterfaceCreate(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := mix.NewDataBag("Nothing Created")
		err := mix.Write(w, fact.Create(r, "Index", "./views/index.html", b))

		if err != nil {
			log.Println(err)
		}
	}
}
