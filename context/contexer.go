package context

import "net/http"

//Contexer provides context around Requests and Responses
type Contexer interface {
	SetHeader(key string, val string)            //SetHeader sets a value on the Response Header
	SetStatus(code int)                          //SetStatus set the final Response Status
	FindParam(name string) string                //FindParam returns the value of a query string parameter
	WriteResponse(data []byte) (int, error)      //WriteResponse writes the data to the ResponseWriter
	RequestURI() string                          //RequestURI returns the full URL Requested
	GetCookie(name string) (*http.Cookie, error) //GetCookie returns the value of a cookie
	Body(container interface{}) error            //Body returns an error when it is unable populate the containercontrollers
}
