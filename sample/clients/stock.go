package clients

import (
	"github.com/louisevanderlith/droxolite/mix"
	"html/template"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

func PartsGet(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Parts", "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsSearch(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Parts", "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsView(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Parts View", "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsCreate(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Parts Create", "./views/stock/parts.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesGet(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Services", "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesSearch(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Services", "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesView(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Services View", "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesCreate(tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage(tmpl, "Services Create", "./views/stock/services.html")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}
