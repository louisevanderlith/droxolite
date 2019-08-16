package xontrols

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/husk"
)

// default paging values
const (
	_page = 1
	_size = 5
)

//APICtrl is used to serve raw HTTP Requests
type APICtrl struct {
	InstanceCtrl
	Data map[string]interface{}
}

//Prepare is called before Invoking the Callback
func (ctrl *APICtrl) Prepare() {
	ctrl.Data = make(map[string]interface{})
	ctrl.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	ctrl.SetHeader("Access-Control-Allow-Credentials", "true")
	ctrl.SetHeader("Server", "kettle")
	ctrl.SetHeader("X-Content-Type-Options", "nosniff")
}

//ServeBinary is used to serve files such as images and documents.
func (ctrl *APICtrl) ServeBinary(data []byte, filename string) {
	dataLen := len(data)
	toTake := 512

	if dataLen < 512 {
		toTake = dataLen
	}

	mimetype := http.DetectContentType(data[:toTake])

	ctrl.ServeBinaryWithMIME(data, filename, mimetype)
}

//ServeBinaryWithMIME is used to serve files such as images and documents. You must specify the MIME Type
func (ctrl *APICtrl) ServeBinaryWithMIME(data []byte, filename, mimetype string) {
	ctrl.SetHeader("Content-Description", "File Transfer")
	//ctrl.SetHeader("Content-Type", mimetype)
	ctrl.SetHeader("Content-Disposition", "attachment; filename="+filename)
	ctrl.SetHeader("Content-Transfer-Encoding", "binary")
	ctrl.SetHeader("Expires", "0")
	ctrl.SetHeader("Cache-Control", "must-revalidate")
	ctrl.SetHeader("Pragma", "public")

	ctrl.ctx.WriteResponse(data)
}

//Serve sends data as JSON response.
func (ctrl *APICtrl) Serve(statuscode int, err error, result interface{}) error {
	resp := bodies.NewRESTResult(statuscode, err, result)

	//ctrl.SetHeader("Content-Type", "application/json; charset=utf-8")

	content, err := json.Marshal(*resp)

	if err != nil {
		ctrl.ctx.SetStatus(http.StatusInternalServerError)
		_, err = ctrl.ctx.WriteResponse([]byte((err.Error())))
		return err
	}

	if resp.Code != http.StatusOK {
		ctrl.ctx.SetStatus(resp.Code)
	}

	_, err = ctrl.ctx.WriteResponse(content)

	return err
}

func (ctrl *APICtrl) Filter(requiredRole roletype.Enum, publicKeyPath, serviceName string) bool {
	path := ctrl.ctx.RequestURI()
	//action := ctrl.ctx.Method()

	if strings.HasPrefix(path, "/favicon") {
		return true
	}

	//requiredRole, err := m.GetRequiredRole(path, action)

	//if err != nil {
	//Missing Mapping, the user doesn't have access to the application
	//	ctx.RenderMethodResult(RenderUnauthorized(err))
	//	return
	//}

	if requiredRole == roletype.Unknown {
		return true
	}

	token, err := getAuthorizationToken(ctrl.ctx)

	if err != nil {
		log.Println(err)
		//ctx.RenderMethodResult(RenderUnauthorized(err))
		return false
	}

	avoc, err := bodies.GetAvoCookie(token, publicKeyPath)

	if err != nil {
		//ctx.RenderMethodResult(RenderUnauthorized(err))
		log.Println(err)
		return false
	}

	allowed, err := bodies.IsAllowed(serviceName, avoc.UserRoles, requiredRole)

	if err != nil {
		//ctx.RenderMethodResult(RenderUnauthorized(err))
		log.Println(err)
		return false
	}

	return allowed
}

//GetKeyedRequest will return the Key and update the Target when Requests are sent for updates.
func (ctrl *APICtrl) GetKeyedRequest(target interface{}) (husk.Key, error) {
	result := struct {
		Key  husk.Key
		Body interface{}
	}{
		Body: target,
	}

	err := ctrl.ctx.Body(&result)

	if err != nil {
		return husk.CrazyKey(), err
	}

	return result.Key, nil
}

//GetPageData turns /B1 into page 1. size 1
func (ctrl *APICtrl) GetPageData() (page, pageSize int) {
	pageData := ctrl.FindParam("pagesize")
	return getPageData(pageData)
}

func getPageData(pageData string) (int, int) {

	if len(pageData) < 2 {
		return _page, _size
	}

	pChar := []rune(pageData[:1])

	if len(pChar) != 1 {
		return _page, _size
	}

	page := int(pChar[0]) % 32
	pageSize, err := strconv.Atoi(pageData[1:])

	if err != nil {
		return _page, _size
	}

	return page, pageSize
}

//Returns the [TOKEN] in 'Bearer [TOKEN]'
func getAuthorizationToken(ctx context.Contexer) (string, error) {
	authHead, err := ctx.GetHeader("Authorization")

	if err != nil {
		return "", err
	}

	parts := strings.Split(authHead, " ")
	tokenType := parts[0]
	if strings.Trim(tokenType, " ") != "Bearer" {
		return "", errors.New("Bearer Authentication only")
	}

	return parts[1], nil
}
