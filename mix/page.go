package mix

import (
	"fmt"
	"github.com/louisevanderlith/droxolite/menu"
	"github.com/louisevanderlith/kong/tokens"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

//Page provides a io.Reader for serving html pages
type pge struct {
	contentPage string
	data        map[string]interface{}
	headers     map[string]string
	template    *template.Template
}

func PreparePage(files *template.Template, name, page string) PageMixer {
	shortName := strings.ToLower(strings.Trim(name, " "))
	htmlName := fmt.Sprintf("%s.html", shortName)

	cpy, err := files.Clone()

	if err != nil {
		panic(err)
	}

	tmpl, err := cpy.ParseFiles(page)

	if err != nil {
		panic(err)
	}

	r := &pge{
		data:     make(map[string]interface{}),
		headers:  make(map[string]string),
		template: tmpl,
	}

	r.contentPage = htmlName
	scriptName := fmt.Sprintf("%s.entry.dart.js", shortName)
	_, err = os.Stat(path.Join("dist/js", scriptName))

	r.data["HasScript"] = err == nil
	r.data["ScriptName"] = scriptName
	r.data["Title"] = name

	return r
}

func (r *pge) ChangeTitle(title string) {
	r.data["Title"] = title
}

func (r *pge) AddMenu(m *menu.Menu) {
	r.data["Menu"] = m
}

func (r *pge) Page(data interface{}, claims tokens.Claimer, token string) Mixer {
	r.data["Data"] = data

	if _, isErr := data.(error); isErr {
		r.data["HasScript"] = false
		r.data["ScriptName"] = ""
		r.contentPage = "error.html"
		return r
	}

	if claims != nil {
		r.data["Identity"] = claims
		r.data["Token"] = token

		//User Details
		if claims.HasUser() {
			//never display the user's key on the front-end
			_, n := claims.GetUserinfo()
			r.data["Username"] = n
		}
	}

	return r
}

func (r *pge) Headers() map[string]string {
	result := make(map[string]string)

	result["X-Frame-Options"] = "SAMEORIGIN"
	result["X-XSS-Protection"] = "1; mode=block"
	result["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	result["Access-Control-Allow-Credentials"] = "true"
	result["Server"] = "kettle"
	result["X-Content-Type-Options"] = "nosniff"

	return result
}

//Reader configures the response for reading
func (r *pge) Reader(w http.ResponseWriter) error {
	return r.template.ExecuteTemplate(w, r.contentPage, r.data)
}
