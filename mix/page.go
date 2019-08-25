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
)

//Page provides a io.Reader for serving html pages
type tmpl struct {
	contentPage string
	Settings    bodies.ThemeSetting
	data        map[string]interface{}
	headers     map[string]string
}

func Page(data interface{}) Mixer {
	result := &tmpl{
		data:    make(map[string]interface{}),
		headers: make(map[string]string),
	}

	if _, isErr := data.(error); isErr {
		result.data["Error"] = data
		result.contentPage = "error.html"
	} else {
		result.data["Data"] = data
	}

	return result
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

	page := r.Settings.Templates.Lookup(contentpage)

	if page == nil {
		return nil, fmt.Errorf("Template not Found: %s", contentpage)
	}

	var buffPage bytes.Buffer
	err := page.ExecuteTemplate(&buffPage, contentpage, r.data)

	if err != nil {
		return nil, err
	}

	r.data["LayoutContent"] = template.HTML(buffPage.String())

	masterPage := r.Settings.Templates.Lookup(r.Settings.MasterTemplate.Name()) // "master.html")
	var buffMaster bytes.Buffer
	err = masterPage.ExecuteTemplate(&buffMaster, r.Settings.MasterTemplate.Name(), r.data)

	if err != nil {
		return nil, err
	}

	return &buffMaster, nil
}

func (r *tmpl) ApplySettings(name string, settings bodies.ThemeSetting, avo *bodies.Cookies) {
	shortName := strings.ToLower(strings.Trim(name, " "))
	if len(r.contentPage) == 0 {
		r.contentPage = fmt.Sprintf("%s.html", shortName)
	}
	r.Settings = settings

	scriptName := fmt.Sprintf("%s.entry.dart.js", shortName)
	_, err := os.Stat(path.Join("dist/js", scriptName))

	r.data["HasScript"] = err == nil
	r.data["ScriptName"] = scriptName

	r.data["ShowSave"] = false
	r.data["Title"] = fmt.Sprintf("%s %s", name, settings.Name)
	r.data["LogoKey"] = settings.LogoKey
	r.data["InstanceID"] = settings.InstanceID
	r.data["Host"] = settings.Host
	r.data["GTag"] = settings.GTag
	r.data["Footer"] = settings.Footer

	//User Details
	loggedIn := avo != nil
	r.data["LoggedIn"] = loggedIn

	if loggedIn {
		r.data["Username"] = avo.Username
	}
}

func (r *tmpl) CreateTopMenu(enablesave bool, menu bodies.Menu) {
	r.data["ShowSave"] = enablesave
	r.data["TopMenu"] = menu.Items()
}

func (r *tmpl) CreateSideMenu(menu *bodies.Menu) {
	r.data["SideMenu"] = menu.Items()
}
