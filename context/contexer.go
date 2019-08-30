package context

import (
	"github.com/louisevanderlith/droxolite/mix"
)

//Contexer provides context around Requests and Responses
type Contexer interface {
	Requester
	Redirect(status int, url string)  //Redirects to the given URL with status code
	SetHeader(key string, val string) //SetHeader sets a value on the Response Header
	SetStatus(code int)               //SetStatus set the final Response Status

	Serve(int, mix.Mixer) error
}
