package models

import (
	"github.com/louisevanderlith/droxolite/security/tracetype"
	"github.com/louisevanderlith/husk"
)

type LoginTrace struct {
	Location string `hsk:"null;size(128)"`
	IP       string `hsk:"null;size(50)"`
	Allowed  bool   `hsk:"default(true)"`
	ClientID string
	TraceEnv tracetype.Enum
}

func (o LoginTrace) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}
