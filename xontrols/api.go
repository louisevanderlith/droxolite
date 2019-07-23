package xontrols

import (
	"encoding/json"
	"net/http"

	"github.com/louisevanderlith/droxolite/bodies"
)

//APICtrl is used to serve raw HTTP Requests
type APICtrl struct {
	InstanceCtrl
	Data map[string]interface{}
}

//Prepare is called before Invoking the Callback
func (ctrl *APICtrl) Prepare() {
	ctrl.Ctx.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	ctrl.Ctx.SetHeader("Access-Control-Allow-Credentials", "true")
	ctrl.Ctx.SetHeader("Server", "kettle")
	ctrl.Ctx.SetHeader("X-Content-Type-Options", "nosniff")
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
	ctrl.Ctx.SetHeader("Content-Description", "File Transfer")
	ctrl.Ctx.SetHeader("Content-Type", mimetype)
	ctrl.Ctx.SetHeader("Content-Disposition", "attachment; filename="+filename)
	ctrl.Ctx.SetHeader("Content-Transfer-Encoding", "binary")
	ctrl.Ctx.SetHeader("Expires", "0")
	ctrl.Ctx.SetHeader("Cache-Control", "must-revalidate")
	ctrl.Ctx.SetHeader("Pragma", "public")

	ctrl.Ctx.WriteResponse(data)
	//Write body(data)
}

//Serve sends data as JSON response.
func (ctrl *APICtrl) Serve(statuscode int, err error, result interface{}) error {
	resp := bodies.NewRESTResult(statuscode, err, result)

	ctrl.Ctx.SetStatus(resp.Code)

	ctrl.Ctx.SetHeader("Content-Type", "application/json; charset=utf-8")

	content, err := json.Marshal(*resp)

	if err != nil {
		ctrl.Ctx.SetStatus(http.StatusInternalServerError)
		_, err = ctrl.Ctx.WriteResponse([]byte((err.Error())))
		return err
	}

	_, err = ctrl.Ctx.WriteResponse(content)

	return err
}

func (ctrl *APICtrl) Filter() bool {
	return true
}
