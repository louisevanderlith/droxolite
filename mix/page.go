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
type tmpl struct {
	contentPage string
	data        map[string]interface{}
	headers     map[string]string
	templates   *template.Template
	master      *template.Template
}

func Page(name string, data interface{}, claims tokens.Claimer, mastr *template.Template, templates *template.Template) Mixer {
	r := &tmpl{
		data:      make(map[string]interface{}),
		headers:   make(map[string]string),
		master:    mastr,
		templates: templates,
	}

	if _, isErr := data.(error); isErr {
		r.data["Error"] = data
		r.contentPage = "error.html"
	} else {
		r.data["Data"] = data
	}

	shortName := strings.ToLower(strings.Trim(name, " "))

	if len(r.contentPage) == 0 {
		r.contentPage = fmt.Sprintf("%s.html", shortName)
	}

	scriptName := fmt.Sprintf("%s.entry.dart.js", shortName)
	_, err := os.Stat(path.Join("dist/js", scriptName))

	r.data["HasScript"] = err == nil
	r.data["ScriptName"] = scriptName
	r.data["Name"] = name

	if claims != nil {
		r.data["Identity"] = claims

		//User Details
		if claims.HasUser() {
			_, n := claims.GetUserinfo()
			r.data["Username"] = n
		}
	}

	return r
}

func (r *tmpl) Headers() map[string]string {
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
func (r *tmpl) Reader() (io.Reader, error) {
	contentpage := r.contentPage

	page := r.templates.Lookup(contentpage)

	if page == nil {
		return nil, fmt.Errorf("Template not Found: %s", contentpage)
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
