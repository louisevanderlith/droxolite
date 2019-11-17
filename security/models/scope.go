package models

import (
	"github.com/louisevanderlith/husk"
)

type Scope struct {
	Enabled     bool
	Name        string `hsk:"size(50)"`
	DisplayName string `hsk:"size(50)"`
	Description string `hsk:"size(256)"`
	Required    bool   //can be disabled on consent screen
	Claims      []Claim
}

func NewScope(name, displayname, description string, claims []Claim) Scope {
	return Scope{
		Enabled:     true,
		Name:        name,
		DisplayName: displayname,
		Description: description,
		Required:    false,
		Claims:      claims,
	}
}

func (s Scope) Valid() (bool, error) {
	return husk.ValidateStruct(&s)
}
