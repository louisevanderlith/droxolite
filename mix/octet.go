package mix

import (
	"bytes"
	"io"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
)

//Octet provides a io.Reader for serving data (octet-)streams
type octet struct {
	filename string
	mimetype string
	headers  map[string]string
	data     []byte
}

func Octet(data interface{}) Mixer {
	result := &octet{
		data: data.([]byte),
	}

	return result
}

//Instead of assigning headers, returns headers that should be applied.
func (r *octet) Headers() map[string]string {
	result := make(map[string]string)

	result["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	result["Access-Control-Allow-Credentials"] = "true"
	result["Server"] = "kettle"
	result["X-Content-Type-Options"] = "nosniff"

	result["Content-Description"] = "File Transfer"
	result["Content-Transfer-Encoding"] = "binary"
	result["Expires"] = "0"
	result["Cache-Control"] = "must-revalidate"
	result["Pragma"] = "public"

	result["Content-Disposition"] = "attachment; filename=" + r.filename
	result["Content-Type"] = r.mimetype

	return result
}

//Reader configures the response for reading
func (r *octet) Reader() (io.Reader, error) {
	return bytes.NewBuffer(r.data), nil
}

func (r *octet) ApplySettings(name string, settings bodies.ThemeSetting) {
	ext := getExt(name)

	mimes := make(map[string]string)
	mimes["js"] = "text/javascript"
	mimes["css"] = "text/css"
	mimes["html"] = "text/html"
	mimes["ico"] = "image/x-icon"
	mimes["font"] = "font/" + ext
	mimes["jpeg"] = "image/jpeg"
	mimes["jpg"] = "image/jpeg"
	mimes["png"] = "image/png"

	r.mimetype = mimes[ext]
	r.filename = name
}

func getExt(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	return filename[dotIndex+1:]
}

/*
//ServeBinaryWithMIME is used to serve files such as images and documents. You must specify the MIME Type
func (ctx *Ctx) ServeBinaryStream(data io.Reader, filename, mimetype string) error {

}

//ServeBinary is used to serve files such as images and documents.
func (ctx *Ctx) ServeBinary(data []byte, filename string) error {
	dataLen := len(data)
	toTake := 512

	if dataLen < 512 {
		toTake = dataLen
	}

	mimetype := http.DetectContentType(data[:toTake])

	return ctx.ServeBinaryWithMIME(data, filename, mimetype)
}

//ServeBinaryWithMIME is used to serve files such as images and documents. You must specify the MIME Type
func (ctx *Ctx) ServeBinaryWithMIME(data []byte, filename, mimetype string) error {
	ctx.SetHeader("Content-Description", "File Transfer")
	ctx.SetHeader("Content-Disposition", "attachment; filename="+filename)
	ctx.SetHeader("Content-Transfer-Encoding", "binary")
	ctx.SetHeader("Expires", "0")
	ctx.SetHeader("Cache-Control", "must-revalidate")
	ctx.SetHeader("Pragma", "public")

	ctx.SetHeader("Content-Type", mimetype)

	_, err := ctx.WriteResponse(data)

	return err
}
*/
