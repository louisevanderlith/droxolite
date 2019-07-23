package context

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Ctx provides context around Requests and Responses
type Ctx struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	//Body           []byte
}

func New(response http.ResponseWriter, request *http.Request) Contexer {
	return &Ctx{
		ResponseWriter: response,
		Request:        request,
	}
}

//SetHeader sets a value on the Response Header
func (ctx *Ctx) SetHeader(key string, val string) {
	ctx.ResponseWriter.Header().Set(key, val)
}

//SetStatus set the final Response Status
func (ctx *Ctx) SetStatus(code int) {
	ctx.ResponseWriter.WriteHeader(code)
}

func (ctx *Ctx) FindQueryParam(name string) string {
	results, ok := ctx.Request.URL.Query()[name]

	if !ok {
		return ""
	}

	return results[0]
}

func (ctx *Ctx) FindParam(name string) string {
	vars := mux.Vars(ctx.Request)

	result, ok := vars[name]

	if !ok {
		return ""
	}

	return result
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

func (ctx *Ctx) Body(container interface{}) error {
	decoder := json.NewDecoder(ctx.Request.Body)

	return decoder.Decode(container)
}
