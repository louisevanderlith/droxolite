package element

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/louisevanderlith/husk"
)

//Identity is the basic controls variables accessed by any web-function
type Identity struct {
	LogoKey        husk.Key
	Name           string //Profile Name 
	Host           string
	InstanceID     string
	GTag           string
	Footer         Footer
	MasterTemplate *template.Template //only has access to _shared
	Templates      *template.Template //has access to _shared and views
}
 
type Footer struct {
	SocialLinks map[string]string //fa-[facebook] = 'www.facebook.com/companyx'
	OtherLinks  map[string]string //fa-[same]
}

func NewIdentity(name, host string, logoKey husk.Key, instanceID, gtag string) *Identity {
	return &Identity{
		Name:       name,
		LogoKey:    logoKey,
		Host:       host,
		InstanceID: instanceID,
		GTag:       gtag,
		Templates:  &template.Template{},
	}
}

func (t *Identity) LoadTemplate(viewPath, masterpage string) error {
	temps, err := template.ParseFiles(findFiles(viewPath)...)

	if err != nil {
		return err
	}

	t.MasterTemplate = template.New(masterpage)
	t.Templates = temps

	return nil
}

func findFiles(templatesPath string) []string {
	var result []string

	filepath.Walk(templatesPath, func(path string, f os.FileInfo, err error) error {

		if err != nil {
			log.Println(err)
		}

		if !f.IsDir() && strings.HasSuffix(path, ".html") {
			log.Println("file", f.Name(), path)
			result = append(result, path)
		}

		return nil
	})

	return result
}
