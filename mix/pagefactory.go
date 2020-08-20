package mix

import (
	"bytes"
	"fmt"
	"github.com/louisevanderlith/droxolite/drx"
	"github.com/louisevanderlith/droxolite/menu"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

type MixerFactory interface {
	ChangeTitle(title string)
	AddMenu(menu *menu.Menu)
	Create(r *http.Request, data interface{}) Mixer
}

func PreparePage(title string, files *template.Template, tmplPath string) MixerFactory {
	cpy, err := files.Clone()

	if err != nil {
		panic(err)
	}

	tmpl, err := cpy.ParseFiles(tmplPath)

	if err != nil {
		panic(err)
	}

	result := &pgeFactory{
		title:    title,
		name:     strings.ToLower(strings.Replace(title, " ", "", -1)),
		template: tmpl,
		model:    make(map[string]interface{}),
	}

	scriptName := fmt.Sprintf("%s.entry.dart.js", result.name)
		_, err = os.Stat(path.Join("dist/js", scriptName))

	result.model["HasScript"] = err == nil
	result.model["ScriptName"] = scriptName

	return result
}

type pgeFactory struct {
	title    string
	name     string
	template *template.Template
	model    map[string]interface{}
}

func (f *pgeFactory) Create(r *http.Request, data interface{}) Mixer {
	f.model["Data"] = data

	claims := drx.GetIdentity(r)
	f.model["Identity"] = claims

	if claims != nil {
		if !strings.Contains(f.title, " - ") {
			profTitle := fmt.Sprintf("%s - %s", f.title, claims.GetProfile())
			f.ChangeTitle(profTitle)
		}

		f.model["Token"] = drx.GetToken(r)

		//User Details
		if claims.HasUser() {
			//never display the user's key on the front-end
			f.model["Username"] = drx.GetUserIdentity(r).GetDisplayName()
		}
	}

	pageBuff := bytes.Buffer{}
	htmlName := fmt.Sprintf("%s.html", f.name)
	err := f.template.ExecuteTemplate(&pageBuff, htmlName, f.model)

	if err != nil {
		panic(err)
	}

	return &pge{
		data:        pageBuff,
		contentPage: f.name,
	}
}

func (f *pgeFactory) ChangeTitle(title string) {
	f.model["Title"] = title
}

func (f *pgeFactory) AddMenu(m *menu.Menu) {
	f.model["Menu"] = m
}
