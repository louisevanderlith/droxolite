package mix

import (
	"bytes"
	"io"
)

//Page provides a io.Reader for serving html pages
type pge struct {
	data        bytes.Buffer
	contentPage string
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
	return &r.data
}
