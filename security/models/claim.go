package models

import (
	"github.com/louisevanderlith/husk"
)

type Claim struct {
	Name        string `hsk:"size(50)"`
	Description string `hsk:"size(256)"`
}

func (c Claim) Valid() (bool, error) {
	return husk.ValidateStruct(&c)
}