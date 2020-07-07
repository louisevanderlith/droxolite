package mix

import (
	"github.com/louisevanderlith/droxolite/menu"
	"github.com/louisevanderlith/kong/tokens"
	"net/http"
)

//Mixer is used by the Contexer to ApplyHeaders and Write the Response from the Reader
type Mixer interface {
	Reader(w http.ResponseWriter) error
	Headers() map[string]string
}

type PageMixer interface {
	Mixer
	ChangeTitle(title string)
	AddMenu(menu *menu.Menu)
	Page(data interface{}, claims tokens.Claimer, token string) Mixer
}

//DefaultHeaders returns a set of Headers that apply to all mixers
func DefaultHeaders() map[string]string {
	result := make(map[string]string)

	result["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	result["Access-Control-Allow-Credentials"] = "true"
	result["Server"] = "kettle"
	result["X-Content-Type-Options"] = "nosniff"

	return result
}
