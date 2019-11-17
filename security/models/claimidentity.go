package models

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type ClaimIdentity struct {
	Subject    string    `json:"sub"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IP         string    `json:"ip"`
	Location   string    `json:"loc"`
	Issuer     string    `json:"iss"`
	Audience   string    `json:"aud"`
	Expiration time.Time `json:"exp"`
	IssuedAt   time.Time `json:"iat"`
	Gravatar   string    `json:"grav"`
	Roles      []Role    `json:"roles"`
	Active     bool      `json:"act"`
	ClientID   string    `json:"client_id"`
	Scope      string    `json:"scope"`
}

func AssembleClaimIdentity(pubKey *rsa.PublicKey, token string) ClaimIdentity {
	return ClaimIdentity{}
}

func (i ClaimIdentity) Claims() jwt.MapClaims {
	result := make(jwt.MapClaims)

	data, err := json.Marshal(i)

	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &result)

	if err != nil {
		return nil
	}

	return result
}

func (i ClaimIdentity) ToJWT(privateKey *rsa.PrivateKey) (string, error) {
	alg := jwt.GetSigningMethod("RS256")

	if alg == nil {
		return "", fmt.Errorf("Couldn't find signing method: %v", "RS256")
	}

	token := jwt.NewWithClaims(alg, i.Claims())

	return token.SignedString(privateKey)
}
