package mix

import (
	"io"
)

//Mixer is used by the drx.ite to ApplyHeaders and Write the Response from the Reader
type Mixer interface {
	Reader() io.Reader
	Headers() map[string]string
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
