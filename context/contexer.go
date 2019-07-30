package context

import (
	"mime/multipart"
	"net/http"
)

//Contexer provides context around Requests and Responses
type Contexer interface {
	Redirect(status int, url string)                                 //Redirects to the given URL with status code
	Scheme() string                                                  //Scheme returns the Input Scheme
	Method() string                                                  //Method returns the Method associated with the Request
	GetHeader(key string) (string, error)                            //GetHeader returns the value of the Request Header
	SetHeader(key string, val string)                                //SetHeader sets a value on the Response Header
	SetStatus(code int)                                              //SetStatus set the final Response Status
	FindParam(name string) string                                    //FindParam returns the value of a path parameter
	FindQueryParam(name string) string                               //FindParam returns the value of a query string parameter
	WriteResponse(data []byte) (int, error)                          //WriteResponse writes the data to the ResponseWriter
	RequestURI() string                                              //RequestURI returns the full URL Requested
	GetCookie(name string) (*http.Cookie, error)                     //GetCookie returns the value of a cookie
	Body(container interface{}) error                                //Body returns an error when it is unable populate the containercontrollers
	Host() string                                                    //Host returns the Hostname of the request
	File(name string) (multipart.File, *multipart.FileHeader, error) //File returns the Request's FileBody, the key should match name
}
