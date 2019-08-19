package context

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

//Ctx provides context around Requests and Responses
type Ctx struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

func New(response http.ResponseWriter, request *http.Request) Contexer {
	return &Ctx{
		ResponseWriter: response,
		Request:        request,
	}
}

//Method returns the Requests' Method
func (ctx *Ctx) Method() string {
	return ctx.Request.Method
}

//GetHeader returns a Request Header
func (ctx *Ctx) GetHeader(key string) (string, error) {
	headers := ctx.Request.Header[key]

	if len(headers) == 0 {
		return "", fmt.Errorf("no header '%s' found", key)
	}

	return headers[0], nil
}

//SetHeader sets a value on the Response Header
func (ctx *Ctx) SetHeader(key string, val string) {
	ctx.ResponseWriter.Header().Set(key, val)
}

//SetStatus set the final Response Status
func (ctx *Ctx) SetStatus(code int) {
	ctx.ResponseWriter.WriteHeader(code)
}

//File returns the Uploaded file.
func (ctx *Ctx) File(name string) (multipart.File, *multipart.FileHeader, error) {
	err := ctx.Request.ParseMultipartForm(32 << 20)

	if err != nil {
		return nil, nil, err
	}

	return ctx.Request.FormFile(name)
}

//FindFormValue is used to read additional information from File Uploads
func (ctx *Ctx) FindFormValue(name string) string {
	return ctx.Request.FormValue(name)
}

//FindQueryParam returns the requested querystring parameter
func (ctx *Ctx) FindQueryParam(name string) string {
	results, ok := ctx.Request.URL.Query()[name]

	if !ok {
		return ""
	}

	return results[0]
}

//FindParam returns the requested path variable
func (ctx *Ctx) FindParam(name string) string {
	vars := mux.Vars(ctx.Request)

	result, ok := vars[name]

	if !ok {
		return ""
	}

	return result
}

func (ctx *Ctx) Redirect(status int, url string) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, url, status)
}

func (ctx *Ctx) WriteResponse(data []byte) (int, error) {
	return ctx.ResponseWriter.Write(data)
}

func (ctx *Ctx) RequestURI() string {
	return ctx.Request.URL.RequestURI()
}

func (ctx *Ctx) GetCookie(name string) (*http.Cookie, error) {
	return ctx.Request.Cookie(name)
}

func (ctx *Ctx) Scheme() string {
	return ctx.Request.URL.Scheme
}

func (ctx *Ctx) Host() string {
	return ctx.Request.Host
}

//Body returns an error when unable to Decode the JSON request Body
func (ctx *Ctx) Body(container interface{}) error {
	decoder := json.NewDecoder(ctx.Request.Body)

	return decoder.Decode(container)
}
