package models

import "time"

type TokenResponse struct {
	Type   string        `json:"typ"` //Beaer
	Expiry time.Duration `json:"exp"`
	Value  string        `json:"val"`
}
