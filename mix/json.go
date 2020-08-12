package mix

import (
	"bytes"
	"encoding/json"
	"io"
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
func (r *js) Reader() io.Reader {
	j, err := json.Marshal(r.data)

	if err != nil {
		panic(err)
		return nil
	}

	return bytes.NewReader(j)
}
