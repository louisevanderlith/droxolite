package mix

import (
	"bytes"
	"html/template"
	"io"
)

//Page provides a io.Reader for serving html pages
type pge struct {
	template    *template.Template
	title       string
	name        string
	model       map[string]interface{}
}

func (r *pge) Headers() map[string]string {
	result := make(map[string]string)

	result["X-Frame-Options"] = "SAMEORIGIN"
	result["X-XSS-Protection"] = "1; mode=block"
	result["Strict-Transport-Security"] = "max-age=31536000; includeSubDomains"
	result["Access-Control-Allow-Credentials"] = "true"
	result["Server"] = "kettle"
	result["X-Content-Type-Options"] = "nosniff"

	return result
}

//Reader configures the response for reading
func (r *pge) Reader() io.Reader {

	pageBuff := bytes.Buffer{}
	err := r.template.ExecuteTemplate(&pageBuff, r.name, r.model)

	if err != nil {
		panic(err)
	}

	return &pageBuff
}
