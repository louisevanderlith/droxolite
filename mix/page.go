package mix

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/element"
)

//Page provides a io.Reader for serving html pages
type tmpl struct {
	contentPage string
	Identity    *element.Identity
	data        map[string]interface{}
	headers     map[string]string
}

func Page(name string, data interface{}, d *element.Identity, avoc *bodies.Cookies) Mixer {
	r := &tmpl{
		data:    make(map[string]interface{}),
		headers: make(map[string]string),
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

	r.Identity = d

	scriptName := fmt.Sprintf("%s.entry.dart.js", shortName)
	_, err := os.Stat(path.Join("dist/js", scriptName))

	r.data["ShowSearch"] = !(strings.HasSuffix(shortName, "create") || strings.HasSuffix(shortName, "view"))
	r.data["HasScript"] = err == nil
	r.data["ScriptName"] = scriptName
	r.data["Name"] = name
	r.data["Identity"] = d

	//User Details
	loggedIn := avoc != nil
	r.data["LoggedIn"] = loggedIn

	if loggedIn {
		r.data["Username"] = avoc.Username
		r.data["Gravatar"] = avoc.Gravatar
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

	page := r.Identity.Templates.Lookup(contentpage)

	if page == nil {
		return nil, fmt.Errorf("Template not Found: %s", contentpage)
	}

	var buffPage bytes.Buffer
	err := page.ExecuteTemplate(&buffPage, contentpage, r.data)

	if err != nil {
		return nil, err
	}

	r.data["LayoutContent"] = template.HTML(buffPage.String())

	masterPage := r.Identity.Templates.Lookup(r.Identity.MasterTemplate.Name())
	var buffMaster bytes.Buffer
	err = masterPage.ExecuteTemplate(&buffMaster, r.Identity.MasterTemplate.Name(), r.data)

	if err != nil {
		return nil, err
	}

	return &buffMaster, nil
}

func (r *tmpl) CreateSideMenu(menu *bodies.Menu) {
	r.data["SideMenu"] = menu.Items()
}
