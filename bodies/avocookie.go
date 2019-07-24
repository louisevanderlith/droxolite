package bodies

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/louisevanderlith/droxolite/roletype"
	"github.com/louisevanderlith/husk"
	secure "github.com/louisevanderlith/secure/core"
)

//Cookies is our Cookie object.
type Cookies struct {
	UserKey    husk.Key
	Username   string
	UserRoles  secure.ActionMap
	IP         string
	Location   string
	Issuer     string    `json:"iss"`
	Audience   string    `json:"aud"`
	Expiration time.Time `json:"exp"`
	IssuedAt   time.Time `json:"iat"`
}

//NewCookies returns some new Cookies.
func NewCookies(userkey husk.Key, username, ip, location string, roles secure.ActionMap) *Cookies {
	return &Cookies{
		UserKey:    userkey,
		Username:   username,
		IP:         ip,
		Location:   location,
		UserRoles:  roles,
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(time.Hour * 6),
		Issuer:     "https://secure.localhost/oauth/",
		Audience:   "https://localhost",
	}
}

//GetClaims return the JWT Claims from the Cookies Object
func (c Cookies) GetClaims() jwt.MapClaims {
	result := make(jwt.MapClaims)

	data, err := json.Marshal(c)

	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &result)

	if err != nil {
		return nil
	}

	return result
}

func GetAvoCookie(sessionID, publickeyPath string) (*Cookies, error) {
	if len(sessionID) == 0 {
		return nil, errors.New("SessionID empty")
	}

	token, err := jwt.Parse(sessionID, func(t *jwt.Token) (interface{}, error) {
		var rdr io.Reader
		if f, err := os.Open(publickeyPath); err == nil {
			rdr = f
			defer f.Close()
		} else {
			return "", err
		}

		bits, err := ioutil.ReadAll(rdr)

		if err != nil {
			return "", err
		}

		return jwt.ParseRSAPublicKeyFromPEM(bits)
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token invalid")
	}

	jClaim, err := json.Marshal(token.Claims)

	if err != nil {
		return nil, err
	}

	result := &Cookies{}
	err = json.Unmarshal(jClaim, result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func IsAllowed(appName string, usrRoles secure.ActionMap, required roletype.Enum) (bool, error) {
	if required == roletype.Unknown {
		return true, nil
	}
	return hasRole(appName, usrRoles, required)
}

func hasRole(appName string, usrRoles secure.ActionMap, required roletype.Enum) (bool, error) {
	role, err := getRole(appName, usrRoles)

	if err != nil {
		return false, err
	}

	return role <= required, nil
}

func getRole(appName string, usrRoles secure.ActionMap) (roletype.Enum, error) {
	role, ok := usrRoles[appName]

	if !ok {
		msg := fmt.Errorf("application permission required. %s", appName)
		return roletype.Unknown, msg
	}

	return role, nil
}

func removeToken(url string) (string, string) {
	idx := strings.LastIndex(url, "?access_token")

	if idx == -1 {
		return url, ""
	}

	tokenIdx := strings.LastIndex(url, "=") + 1

	cleanURL := url[:idx]
	token := url[tokenIdx:]

	return cleanURL, token
}
