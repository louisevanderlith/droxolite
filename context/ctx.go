package context

import (
	"encoding/json"
	"fmt"
	"github.com/louisevanderlith/kong/tokens"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/husk"
)

//Ctx provides context around Requests and Responses
type Ctx struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	claims         tokens.Claimer
}

func New(w http.ResponseWriter, r *http.Request) Contexer {
	result := &Ctx{
		responseWriter: w,
		request:        r,
	}

	clms := r.Context().Value("claims")
	if clms != nil {
		result.claims = clms.(tokens.Claimer)
	}

	return result
}

func (ctx *Ctx) Request() *http.Request {
	return ctx.request
}

func (ctx *Ctx) Responder() http.ResponseWriter {
	return ctx.responseWriter
}

//Method returns the Requests' Method
func (ctx *Ctx) Method() string {
	return ctx.request.Method
}

//GetHeader returns a Request Header
func (ctx *Ctx) GetHeader(key string) (string, error) {
	headers := ctx.request.Header[key]

	if len(headers) == 0 {
		return "", fmt.Errorf("no header '%s' found", key)
	}

	return headers[0], nil
}

//SetHeader sets a value on the Response Header
func (ctx *Ctx) SetHeader(key string, val string) {
	ctx.responseWriter.Header().Set(key, val)
}

//SetStatus set the final Response Status
func (ctx *Ctx) SetStatus(code int) {
	ctx.responseWriter.WriteHeader(code)
}

//File returns the Uploaded file.
func (ctx *Ctx) File(name string) (multipart.File, *multipart.FileHeader, error) {
	err := ctx.request.ParseMultipartForm(32 << 20)

	if err != nil {
		return nil, nil, err
	}

	return ctx.request.FormFile(name)
}

//FindFormValue is used to read additional information from File Uploads
func (ctx *Ctx) FindFormValue(name string) string {
	return ctx.request.FormValue(name)
}

//FindQueryParam returns the requested querystring parameter
func (ctx *Ctx) FindQueryParam(name string) string {
	results, ok := ctx.request.URL.Query()[name]

	if !ok {
		return ""
	}

	return results[0]
}

//FindParam returns the requested path variable
func (ctx *Ctx) FindParam(name string) string {
	vars := mux.Vars(ctx.request)

	result, ok := vars[name]

	if !ok {
		return ""
	}

	return result
}

func (ctx *Ctx) Redirect(status int, url string) {
	http.Redirect(ctx.responseWriter, ctx.request, url, status)
}

func (ctx *Ctx) WriteResponse(data []byte) (int, error) {
	return ctx.responseWriter.Write(data)
}

func (ctx *Ctx) WriteStreamResponse(data io.Reader) (int64, error) {
	return io.Copy(ctx.responseWriter, data)
}

func (ctx *Ctx) RequestURI() string {
	return ctx.request.URL.RequestURI()
}

func (ctx *Ctx) GetCookie(name string) (*http.Cookie, error) {
	return ctx.request.Cookie(name)
}

func (ctx *Ctx) Scheme() string {
	return ctx.request.URL.Scheme
}

func (ctx *Ctx) Host() string {
	return ctx.request.Host
}

//Body returns an error when unable to Decode the JSON request
func (ctx *Ctx) Body(container interface{}) error {
	decoder := json.NewDecoder(ctx.request.Body)

	return decoder.Decode(container)
}

func (ctx *Ctx) GetInstanceID() string {
	return ctx.claims.GetId()
}

//Serve is usually sent a Mixer. Serve(mixer.JSON(500, nil))
func (ctx *Ctx) Serve(status int, mx mix.Mixer) error {
	for key, head := range mx.Headers() {
		ctx.SetHeader(key, head)
	}

	if status != http.StatusOK {
		ctx.SetStatus(status)
	}

	return mx.Reader(ctx.responseWriter)
}

//GetKeyedRequest will return the Key and update the Target when Requests are sent for updates.
func (ctx *Ctx) GetKeyedRequest(target interface{}) (husk.Key, error) {
	result := struct {
		Key  husk.Key
		Body interface{}
	}{
		Body: target,
	}

	err := ctx.Body(&result)

	if err != nil {
		return husk.CrazyKey(), err
	}

	return result.Key, nil
}

//GetPageData turns /B1 into page 1. size 1
func (ctx *Ctx) GetPageData() (page, pageSize int) {
	pageData := ctx.FindParam("pagesize")
	return getPageData(pageData)
}

func getPageData(pageData string) (int, int) {
	defaultPage := 1
	defaultSize := 10

	if len(pageData) < 2 {
		return defaultPage, defaultSize
	}

	pChar := []rune(pageData[:1])

	if len(pChar) != 1 {
		return defaultPage, defaultSize
	}

	page := int(pChar[0]) % 32
	pageSize, err := strconv.Atoi(pageData[1:])

	if err != nil {
		return defaultPage, defaultSize
	}

	return page, pageSize
}

func (ctx *Ctx) GetToken() string {
	v := ctx.Request().Context().Value("token")

	tkn, ok := v.(string)

	if !ok {
		return ""
	}

	return tkn
}

func (ctx *Ctx) GetTokenInfo() tokens.Claimer {
	return ctx.claims
}
