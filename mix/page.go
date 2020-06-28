package mix

import (
	"bytes"
	"fmt"
	"github.com/louisevanderlith/kong/tokens"
	"html/template"
	"io"
	"os"
	"path"
	"strings"
)

//Page provides a io.Reader for serving html pages
type pge struct {
	contentPage string
	data        map[string]interface{}
	headers     map[string]string
	templates   *template.Template
	master      *template.Template
}

func PreparePage(title, name string, mastr *template.Template, templates *template.Template) PageMixer {
	r := &pge{
		data:      make(map[string]interface{}),
		headers:   make(map[string]string),
		master:    mastr,
		templates: templates,
	}

	shortName := strings.ToLower(strings.Trim(name, " "))
	r.contentPage = fmt.Sprintf("%s.html", shortName)

	scriptName := fmt.Sprintf("%s.entry.dart.js", shortName)
	_, err := os.Stat(path.Join("dist/js", scriptName))

	r.data["HasScript"] = err == nil
	r.data["ScriptName"] = scriptName
	r.data["Title"] = title

	return r
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
func (r *pge) Reader() (io.Reader, error) {
	contentpage := r.contentPage

	page := r.templates.Lookup(contentpage)

	if page == nil {
		return nil, fmt.Errorf("template not found: %s", contentpage)
	}

	var buffPage bytes.Buffer
	err := page.ExecuteTemplate(&buffPage, contentpage, r.data)

	if err != nil {
		return nil, err
	}

	r.data["LayoutContent"] = template.HTML(buffPage.String())

	masterPage := r.templates.Lookup(r.master.Name())
	var buffMaster bytes.Buffer
	err = masterPage.ExecuteTemplate(&buffMaster, r.master.Name(), r.data)

	if err != nil {
		return nil, err
	}

	return &buffMaster, nil
}
