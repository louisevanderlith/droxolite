package mix

import (
	"encoding/json"
	"net/http"
)

//JSON provides a io.Reader for serving json data
type js struct {
	headers map[string]string
	data    interface{}
}

//JSON is called before every function execution to setup the environment a Handler will expect
func JSON(data interface{}) Mixer {
	result := &js{
		headers: DefaultHeaders(),
		data:    data,
	}

	return result
}

func (r *js) Headers() map[string]string {
	return r.headers
}

//Reader configures the response for reading
func (r *js) Reader(w http.ResponseWriter) error {
	enc := json.NewEncoder(w)
	return enc.Encode(r.data)
}
