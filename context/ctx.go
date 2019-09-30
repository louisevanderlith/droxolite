package context

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/mix"
	"github.com/louisevanderlith/husk"
)

//Ctx provides context around Requests and Responses
type Ctx struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	instanceID     string
	publicKey      string
}

func New(response http.ResponseWriter, request *http.Request, instanceID, publicKey string) Contexer {
	return &Ctx{
		ResponseWriter: response,
		Request:        request,
		instanceID:     instanceID,
		publicKey:      publicKey,
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

func (ctx *Ctx) WriteStreamResponse(data io.Reader) (int64, error) {
	return io.Copy(ctx.ResponseWriter, data)
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

//Body returns an error when unable to Decode the JSON request
func (ctx *Ctx) Body(container interface{}) error {
	decoder := json.NewDecoder(ctx.Request.Body)

	return decoder.Decode(container)
}

func (ctx *Ctx) GetInstanceID() string {
	return ctx.instanceID
}

//Serve is usually sent a Mixer. Serve(mixer.JSON(500, nil))
func (ctx *Ctx) Serve(status int, mx mix.Mixer) error {
	for key, head := range mx.Headers() {
		ctx.SetHeader(key, head)
	}

	if status != http.StatusOK {
		ctx.SetStatus(status)
	}

	readr, err := mx.Reader()

	if err != nil {
		return err
	}

	_, err = io.Copy(ctx.ResponseWriter, readr)

	return err
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

func (ctx *Ctx) GetMyToken() string {
	cooki, err := ctx.GetCookie("avosession")

	if err != nil {
		return ""
	}

	return cooki.Value
}

func (ctx *Ctx) GetMyUser() *bodies.Cookies {
	token := ctx.GetMyToken()

	avoc, err := bodies.GetAvoCookie(token, ctx.publicKey)

	if err != nil {
		log.Println(err)
		return nil
	}

	return avoc
}
