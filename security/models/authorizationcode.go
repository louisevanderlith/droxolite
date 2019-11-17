package models

import (
	"net/http"

	"github.com/louisevanderlith/husk"
)

type Authorization struct {
	ClientID     husk.Key
	Audience     string
	Scope        string
	ResponseType string
	State        string
	RedirectUri  string
}

func (a Authorization) Authorize() (int, interface{}) {
	return http.StatusNotImplemented, nil
}
