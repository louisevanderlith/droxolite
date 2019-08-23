package mix

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
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

	masterPage := r.Settings.Templates.Lookup("master.html")
	var buffMaster bytes.Buffer
	err = masterPage.ExecuteTemplate(&buffMaster, "master.html", r.data)

	if err != nil {
		return nil, err
	}

	return &buffMaster, nil
}

func (ctrl *tmpl) ApplySettings(name string, settings bodies.ThemeSetting) {
	if len(ctrl.contentPage) == 0 {
		ctrl.contentPage = fmt.Sprintf("%s.html", strings.ToLower(strings.Trim(name, " ")))
	}
	ctrl.Settings = settings
}

func (ctrl *tmpl) EnableSave() {
	ctrl.data["ShowSave"] = true
}

/*
func (ctrl *HTML) Setup(name, title string, hasScript bool) {
	ctrl.ContentPage = fmt.Sprintf("%s.html", name)
	ctrl.applySettings(title)

	ctrl.Data["HasScript"] = hasScript
	ctrl.Data["ScriptName"] = fmt.Sprintf("%s.entry.dart.js", name)
	ctrl.Data["ShowSave"] = false
}

func (ctrl *HTML) applySettings(title string) {
	ctrl.Data["Title"] = fmt.Sprintf("%s %s", title, ctrl.Settings.Name)
	ctrl.Data["LogoKey"] = ctrl.Settings.LogoKey
	ctrl.Data["InstanceID"] = ctrl.Settings.InstanceID
	ctrl.Data["Host"] = ctrl.Settings.Host
	ctrl.Data["GTag"] = ctrl.Settings.GTag
	ctrl.Data["Footer"] = ctrl.Settings.Footer

	//User Details
	loggedIn := ctrl.AvoCookie != nil
	ctrl.Data["LoggedIn"] = loggedIn

	if loggedIn {
		ctrl.Data["Username"] = ctrl.AvoCookie.Username
	}
}

//CreateTopMenu sets the content of the Top menu bar
func (ctrl *HTML) CreateTopMenu(enablesave bool, menu bodies.Menu) {
	ctrl.Data["TopMenu"] = menu.Items()
}

func (ctrl *HTML) CreateSideMenu(menu *bodies.Menu) {
	ctrl.Data["SideMenu"] = menu.Items()
}
*/
