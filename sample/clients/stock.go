package clients

import (
	"github.com/louisevanderlith/droxolite/mix"
	"html/template"
	"log"
	"net/http"

	"github.com/louisevanderlith/droxolite/context"
)

func PartsGet(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts!", "index", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsSearch(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts!", "parts", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsView(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts!", "parts", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func PartsCreate(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Parts!", "parts", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesGet(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services!", "services", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesSearch(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services!", "services", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesView(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services!", "services", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}

func ServicesCreate(mstr *template.Template, tmpl *template.Template) http.HandlerFunc {
	pge := mix.PreparePage("Services!", "services", mstr, tmpl)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(w, r)
		mxr := pge.Page(nil, ctx.GetTokenInfo(), ctx.GetToken())

		err := ctx.Serve(http.StatusOK, mxr)

		if err != nil {
			log.Println(err)
		}
	}
}
