package context

import (
	"github.com/louisevanderlith/droxolite/mix"
)

//Contexer provides context around Requests and Responses
type Contexer interface {
	Requester
	SetHeader(key string, val string) //SetHeader sets a value on the Response Header
	SetStatus(code int)               //SetStatus set the final Response Status

	Serve(int, mix.Mixer) error
}
