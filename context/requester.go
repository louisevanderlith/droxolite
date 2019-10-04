package context

import (
	"mime/multipart"
	"net/http"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/husk"
)

type Requester interface {
	Request() *http.Request                                          //Request returns the current http request
	Responder() http.ResponseWriter                                  //Responder returns the current http response writer
	GetInstanceID() string                                           //GetInstanceID returns the Service's ID
	Scheme() string                                                  //Scheme returns the Input Scheme
	Method() string                                                  //Method returns the Method associated with the Request
	GetHeader(key string) (string, error)                            //GetHeader returns the value of the Request Header
	FindParam(name string) string                                    //FindParam returns the value of a path parameter
	FindQueryParam(name string) string                               //FindParam returns the value of a query string parameter
	RequestURI() string                                              //RequestURI returns the full URL Requested
	GetCookie(name string) (*http.Cookie, error)                     //GetCookie returns the value of a cookie
	Body(container interface{}) error                                //Body returns an error when it is unable populate the container controllers
	Host() string                                                    //Host returns the Hostname of the request
	File(name string) (multipart.File, *multipart.FileHeader, error) //File returns the Request's FileBody, the key should match name
	FindFormValue(name string) string                                //FindFormValue is used to read additional information from File Uploads
	GetKeyedRequest(target interface{}) (husk.Key, error)
	GetPageData() (page, pageSize int)
	GetMyToken() string
	GetMyUser() *bodies.Cookies
	Redirect(status int, url string) //Redirects to the given URL with status code
}
