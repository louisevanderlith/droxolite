package clients

import (
	"github.com/louisevanderlith/droxolite/mix"
	"html/template"
	"log"
	"net/http"
)

func PartsGet(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts", tmpl, "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsSearch(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts", tmpl, "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsView(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts View", tmpl, "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsCreate(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts Create", tmpl, "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesGet(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services", tmpl, "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesSearch(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services", tmpl, "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesView(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services View", tmpl, "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesCreate(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services Create", tmpl, "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		err := mix.Write(w, pge.Create(r, nil))

		if err != nil {
			log.Println(err)
		}
	}
}
