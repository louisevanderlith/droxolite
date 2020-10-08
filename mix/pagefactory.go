package mix

import (
	"bytes"
	"fmt"
	"github.com/louisevanderlith/droxolite/menu"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

type MixerFactory interface {
	SetValue(name string, val interface{})
	GetTitle() string
	ChangeTitle(title string)
	AddMenu(menu *menu.Menu)
	Create(r *http.Request, data interface{}) Mixer
	AddModifier(mod func(f MixerFactory, r *http.Request))
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
		name:     fetchName(tmplPath),
		template: tmpl,
		model:    make(map[string]interface{}),
	}

	baseName := strings.ToLower(strings.Replace(title, " ", "", -1))
	scriptName := fmt.Sprintf("%s.entry.dart.js", baseName)
	_, err = os.Stat(path.Join("dist/js", scriptName))

	result.model["HasScript"] = err == nil
	result.model["ScriptName"] = scriptName

	return result
}

func fetchName(tmplPath string) string {
	lastSlash := strings.LastIndex(tmplPath, "/")
	return strings.ToLower(tmplPath[(lastSlash + 1):])
}

type pgeFactory struct {
	title     string
	name      string
	template  *template.Template
	model     map[string]interface{}
	modifiers []func(f MixerFactory, r *http.Request)
}

func (f *pgeFactory) GetTitle() string {
	return f.title
}

func (f *pgeFactory) AddModifier(mod func(f MixerFactory, r *http.Request)) {
	f.modifiers = append(f.modifiers, mod)
}

func (f *pgeFactory) Create(r *http.Request, data interface{}) Mixer {
	f.model["Data"] = data

	for _, mod := range f.modifiers {
		mod(f, r)
	}

	pageBuff := bytes.Buffer{}
	err := f.template.ExecuteTemplate(&pageBuff, f.name, f.model)

	if err != nil {
		panic(err)
	}

	return &pge{
		data:        pageBuff,
		contentPage: f.name,
	}
}

func (f *pgeFactory) ChangeTitle(title string) {
	f.SetValue("Title", title)
}

func (f *pgeFactory) AddMenu(m *menu.Menu) {
	f.SetValue("Menu", m)
}

func (f *pgeFactory) SetValue(name string, val interface{}) {
	f.model[name] = val
}
