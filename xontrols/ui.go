package xontrols

import (
	"fmt"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
)

//UICtrl is the base for all APP Controllers
type UICtrl struct {
	APICtrl
	Layout     string
	MasterPage string //Base Template page.
	Settings   bodies.ThemeSetting
}

func (ctrl *UICtrl) SetTheme(settings bodies.ThemeSetting) {
	ctrl.Settings = settings
}

func (ctrl *UICtrl) Prepare() {
	defer ctrl.APICtrl.Prepare()

	ctrl.Layout = "_shared/master.html"
	ctrl.Ctx.SetHeader("X-Frame-Options", "SAMEORIGIN")
	ctrl.Ctx.SetHeader("X-XSS-Protection", "1; mode=block")
}

func (ctrl *UICtrl) Setup(name, title string, hasScript bool) {
	ctrl.MasterPage = fmt.Sprintf("%s.html", name)
	ctrl.applySettings(title)

	ctrl.Data["HasScript"] = hasScript
	ctrl.Data["ScriptName"] = fmt.Sprintf("%s.entry.dart.js", name)
	ctrl.Data["ShowSave"] = false
}

func (ctrl *UICtrl) EnableSave() {
	ctrl.Data["ShowSave"] = true
}

func (ctrl *UICtrl) applySettings(title string) {
	ctrl.Data["Title"] = fmt.Sprintf("%s %s", title, ctrl.Settings.Name)
	ctrl.Data["LogoKey"] = ctrl.Settings.LogoKey
	ctrl.Data["InstanceID"] = ctrl.Settings.InstanceID
	ctrl.Data["Host"] = ctrl.Settings.Host
	ctrl.Data["Crumbs"] = decipherURL(ctrl.Ctx.RequestURI())
	ctrl.Data["GTag"] = ctrl.Settings.GTag

	//User Details

	loggedIn := ctrl.AvoCookie != nil
	ctrl.Data["LoggedIn"] = loggedIn

	if loggedIn {
		ctrl.Data["Username"] = ctrl.AvoCookie.Username
	}
}

//Serve sends the response with 'Error' and 'Data' properties.
func (ctrl *UICtrl) Serve(data interface{}, err error) {
	if err != nil {
		ctrl.Ctx.SetStatus(500)
	}

	ctrl.Data["Error"] = err
	ctrl.Data["Data"] = data
}

func (ctrl *UICtrl) Filter() bool {
	return true
}

//ServeJSON enables JSON Responses on UI Controllers
func (ctrl *UICtrl) ServeJSON(statuscode int, err error, data interface{}) {
	//ctrl.EnableRender = false

	ctrl.APICtrl.Serve(statuscode, err, data)
	//ctrl.EnableRender = true
}

func (ctrl *UICtrl) CreateTopMenu(menu *bodies.Menu) {
	ctrl.Data["TopMenu"] = menu
}

func (ctrl *UICtrl) CreateSideMenu(menu *bodies.Menu) {
	ctrl.Data["SideMenu"] = menu
}

func (ctrl *UICtrl) GetMyToken() string {
	cooki, err := ctrl.Ctx.GetCookie("avosession")

	if err != nil {
		return ""
	}

	return cooki.Value
}

func decipherURL(url string) []string {
	var result []string
	qryIndx := strings.Index(url, "?")

	if qryIndx != -1 {
		url = url[:qryIndx]
	}

	parts := strings.Split(url, "/")

	for _, v := range parts {
		if len(v) > 0 {
			result = append(result, v)
		}
	}

	return result
}
