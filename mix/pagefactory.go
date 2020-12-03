package mix

import (
	"fmt"
	"github.com/louisevanderlith/droxolite/menu"
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

type MixerFactory interface {
	AddMenu(menu *menu.Menu)
	Create(r *http.Request, title, path string, data Bag) Mixer
	AddModifier(mod ModFunc)
}

func NewPageFactory(files *template.Template, mods ...ModFunc) MixerFactory {
	return &pgeFactory{
		files:     files,
		modifiers: mods,
	}
}

func fetchName(tmplPath string) string {
	lastSlash := strings.LastIndex(tmplPath, "/")
	return strings.ToLower(tmplPath[(lastSlash + 1):])
}

type pgeFactory struct {
	files     *template.Template
	modifiers []ModFunc
	menu      *menu.Menu
}

func (f *pgeFactory) AddModifier(mod ModFunc) {
	f.modifiers = append(f.modifiers, mod)
}

func (f *pgeFactory) Create(r *http.Request, title, templatePath string, bag Bag) Mixer {
	if bag == nil {
		bag = NewBag()
	}

	baseName := strings.ToLower(strings.Replace(title, " ", "", -1))
	scriptName := fmt.Sprintf("%s.entry.dart.js", baseName)
	_, err := os.Stat(path.Join("dist/js", scriptName))

	bag.SetValue("HasScript", err == nil)
	bag.SetValue("ScriptName", scriptName)

	if f.menu != nil {
		bag.SetValue("Menu", f.menu)
	}

	for _, mod := range f.modifiers {
		mod(bag, r)
	}

	cpy, err := f.files.Clone()

	if err != nil {
		panic(err)
	}

	tmpl, err := cpy.ParseFiles(templatePath)

	if err != nil {
		panic(err)
	}

	return &pge{
		template: tmpl,
		title:    title,
		name:     fetchName(templatePath),
		model:    bag.Values(),
	}
}

func (f *pgeFactory) AddMenu(m *menu.Menu) {
	f.menu = m
}
