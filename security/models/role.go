package models

import (
	"github.com/louisevanderlith/droxolite/security/roletype"
	"github.com/louisevanderlith/husk"
)

type Role struct {
	ApplicationName string
	Description     roletype.Enum
}

func (o Role) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}
