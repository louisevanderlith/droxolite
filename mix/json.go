package mix

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/louisevanderlith/droxolite/bodies"
)

// default paging values
const (
	_page = 1
	_size = 5
)

//JSON provides a io.Reader for serving json data
type js struct {
	headers map[string]string
	data    interface{}
}

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
func (r *js) Reader() (io.Reader, error) {
	resp := bodies.NewRESTResult(r.data)

	content, err := json.Marshal(*resp)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(content), nil
	//_, err = ctrl.ctx.WriteResponse(content)

	//return err
}

func (r *js) ApplySettings(name string, settings bodies.ThemeSetting, avo *bodies.Cookies) {

}

/*
//Serve sends data as JSON response.
func (ctx *JSON) ServeJSON(statuscode int, err error, result interface{}) error {

}
*/
