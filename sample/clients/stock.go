package clients

import (
	"github.com/louisevanderlith/droxolite/mix"
	"log"
	"net/http"
)

func PartsGet(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, fact.Create(r, "Parts", "./views/stock/parts.html", nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsSearch(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, fact.Create(r, "Parts", "./views/stock/parts.html", nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsView(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, fact.Create(r, "Parts View", "./views/stock/parts.html", nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsCreate(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, fact.Create(r, "Parts Create", "./views/stock/parts.html", nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesGet(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mx := fact.Create(r, "Services", "./views/stock/services.html", nil)
		err := mix.Write(w, mx)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesSearch(fact mix.MixerFactory) http.HandlerFunc {
	//pge := mix.PreparePage("Services", tmpl, "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		mx := fact.Create(r, "Services", "./views/stock/services.html", nil)
		err := mix.Write(w, mx)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesView(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mx := fact.Create(r, "Services", "./views/stock/services.html", nil)
		err := mix.Write(w, mx)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesCreate(fact mix.MixerFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, fact.Create(r, "Services Create", "./views/stock/services.html", nil))

		if err != nil {
			log.Println(err)
		}
	}
}
