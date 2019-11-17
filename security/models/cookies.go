package models

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/louisevanderlith/droxolite/security/roletype"
	"github.com/louisevanderlith/husk"
)

//Cookies is our Claims object
type Cookies struct {
	UserKey    husk.Key
	Username   string
	UserRoles  map[string]int
	IP         string
	Location   string
	Issuer     string    `json:"iss"`
	Audience   string    `json:"aud"`
	Expiration time.Time `json:"exp"`
	IssuedAt   time.Time `json:"iat"`
	Gravatar   string
}

//NewCookies returns some new Claims.
func NewCookies(userkey husk.Key, username, ip, location, email string, roles map[string]int) *Cookies {
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
		Gravatar:   hashGravatar(email),
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

func IsAllowed(appName string, usrRoles map[string]int, required roletype.Enum) (bool, error) {
	if required == roletype.Unknown {
		return true, nil
	}
	return hasRole(appName, usrRoles, required)
}

func hasRole(appName string, usrRoles map[string]int, required roletype.Enum) (bool, error) {
	role, err := getRole(appName, usrRoles)

	if err != nil {
		return false, err
	}

	return role <= required, nil
}

func getRole(appName string, usrRoles map[string]int) (roletype.Enum, error) {
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

func hashGravatar(email string) string {
	h := md5.New()
	io.WriteString(h, strings.ToLower(strings.Trim(email, " ")))

	return fmt.Sprintf("%x", h.Sum(nil))
}
