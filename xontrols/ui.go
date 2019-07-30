package xontrols

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/roletype"
)

//UICtrl is the base for all APP Controllers
type UICtrl struct {
	APICtrl
	MasterPage  string //Base Template page.
	ContentPage string
	Settings    bodies.ThemeSetting
}

func (ctrl *UICtrl) SetTheme(settings bodies.ThemeSetting, masterpage string) {
	ctrl.MasterPage = masterpage
	ctrl.Settings = settings
}

func (ctrl *UICtrl) Prepare() {
	defer ctrl.APICtrl.Prepare()

	ctrl.SetHeader("X-Frame-Options", "SAMEORIGIN")
	ctrl.SetHeader("X-XSS-Protection", "1; mode=block")
}

func (ctrl *UICtrl) Setup(name, title string, hasScript bool) {
	ctrl.ContentPage = fmt.Sprintf("%s.html", name)
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
	ctrl.Data["Crumbs"] = decipherURL(ctrl.ctx.RequestURI())
	ctrl.Data["GTag"] = ctrl.Settings.GTag

	//User Details
	loggedIn := ctrl.AvoCookie != nil
	ctrl.Data["LoggedIn"] = loggedIn

	if loggedIn {
		ctrl.Data["Username"] = ctrl.AvoCookie.Username
	}
}

//Serve sends the response with 'Error' and 'Data' properties.
func (ctrl *UICtrl) Serve(statuscode int, err error, result interface{}) error {
	ctrl.ctx.SetStatus(statuscode)
	renderPage := ctrl.ContentPage

	if err != nil {
		ctrl.Data["Error"] = err
		renderPage = "error.html"
	} else {
		ctrl.Data["Data"] = result
	}

	/*if ctrl.Settings.Templates.Tree != nil {
		return errors.New("tree is empty")
	}*/

	//page, err := renderTemplate(renderPage, ctrl.Data)
	page := ctrl.Settings.Templates.Lookup(renderPage)
	var buffPage bytes.Buffer
	err = page.ExecuteTemplate(&buffPage, renderPage, ctrl.Data)

	if err != nil {
		return err
	}

	//masterPage, err := renderTemplate(ctrl.MasterPage, ctrl.Data)
	ctrl.Data["LayoutContent"] = template.HTML(buffPage.String())
	masterPage := ctrl.Settings.Templates.Lookup(ctrl.MasterPage)
	var buffMaster bytes.Buffer
	err = masterPage.ExecuteTemplate(&buffMaster, ctrl.MasterPage, ctrl.Data)

	if err != nil {
		return err
	}

	_, err = ctrl.ctx.WriteResponse(buffMaster.Bytes())

	return err
}

func (ctrl *UICtrl) Filter(requiredRole roletype.Enum, publicKeyPath, serviceName string) bool {
	path := ctrl.ctx.RequestURI()
	//action := ctrl.ctx.Method()

	if strings.HasPrefix(path, "/static") || strings.HasPrefix(path, "/favicon") {
		return true
	}

	//requiredRole, err := ctrl..GetRequiredRole(path, action)

	//if err != nil {
	//	log.Println(err)
	//Missing Mapping, the user doesn't have access to the application, and must request it.
	//	sendToSubscription(ctx, m.GetInstanceID())
	//	return
	//}

	if requiredRole == roletype.Unknown {
		return true
	}

	token := ctrl.ctx.FindQueryParam("access_token")
	//_, token := removeToken(path)

	if token == "" {
		cookie, err := ctrl.ctx.GetCookie("avosession")

		if err != nil {
			log.Println(err)
			return false
		}

		token = cookie.Value

		if len(token) == 0 {
			//sendToLogin(ctx, ctrl.GetInstanceID())
			return false
		}
	}

	avoc, err := bodies.GetAvoCookie(token, publicKeyPath)

	if err != nil {
		log.Println(err)
		//sendToLogin(ctx, ctrl.GetInstanceID())
		return false
	}

	allowed, err := bodies.IsAllowed(ctrl.Settings.Name, avoc.UserRoles, requiredRole)

	if err != nil || !allowed {
		log.Println(err)
		//sendToLogin(ctx, ctrl.GetInstanceID())
		return false
	}

	return true
}

//ServeJSON enables JSON Responses on UI Controllers
func (ctrl *UICtrl) ServeJSON(statuscode int, err error, data interface{}) {
	ctrl.APICtrl.Serve(statuscode, err, data)
}

//CreateTopMenu sets the content of the Top menu bar
func (ctrl *UICtrl) CreateTopMenu(menu *bodies.Menu) {
	ctrl.Data["TopMenu"] = menu
}

func (ctrl *UICtrl) CreateSideMenu(menu *bodies.Menu) {
	ctrl.Data["SideMenu"] = menu
}

func (ctrl *UICtrl) GetMyToken() string {
	cooki, err := ctrl.ctx.GetCookie("avosession")

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
