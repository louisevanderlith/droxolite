package sample

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/louisevanderlith/droxolite/resins"
)

func GetResponse(e resins.Epoxi, path string, data io.Reader) (*httptest.ResponseRecorder, error) {
	method := "GET"

	if data != nil {
		method = "POST"
	}

	req, err := http.NewRequest(method, path, data)
	req.Host = "localhost"

	if err != nil {
		return nil, err
	}

	handle := e.Router()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	return rr, nil
}
