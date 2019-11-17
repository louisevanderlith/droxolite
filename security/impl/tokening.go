package impl

import (
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	"github.com/louisevanderlith/droxolite/security/client"
	"github.com/louisevanderlith/droxolite/security/tracetype"

	"github.com/louisevanderlith/droxolite/security/models"

	"golang.org/x/crypto/bcrypt"
)

//Token is used to send information to the /token endpoint.
type Tokening struct {
	GrantType    string //Mandatory
	ClientID     string //Mandatory
	ClientSecret string //[Optional]

	//Auth Code Flow
	Code        string
	RedirectURI string

	//Resource Owner Flow
	Audience string
	Scope    string
	Username string
	Password string
}

func (t Tokening) GrantPassword(issuerID string, privateKey *rsa.PrivateKey, store client.UserStorer, registar client.Registar) (int, interface{}) {
	//obj.Username, obj.Password, obj.Audience, obj.Scope, obj.ClientID, obj.ClientSecret
	cred := models.ClientCred{t.ClientID, t.ClientSecret}
	c, err := registar.Authenticate(cred)

	if err != nil {
		return http.StatusUnauthorized, "client is invalid"
	}

	if len(t.Password) < 6 {
		return http.StatusUnauthorized, "user invalid"
	}

	if !strings.Contains(t.Username, "@") {
		return http.StatusUnauthorized, "user invalid"
	}

	subj, user, err := store.FindUser(t.Username)

	if err != nil {
		return http.StatusUnauthorized, err
	}

	if !user.Verified {
		return http.StatusNotAcceptable, "user not yet verified"
	}

	compare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(t.Password))
	passed := compare == nil

	trace := models.LoginTrace{
		Allowed:  passed,
		ClientID: t.ClientID,
		IP:       "none",
		Location: "none",
		TraceEnv: tracetype.Token,
	}

	user.LoginTraces = append(user.LoginTraces, trace)

	err = store.Update(user)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !passed {
		return http.StatusUnauthorized, "login failed"
	}

	exp := time.Hour * 2
	claimid := models.ClaimIdentity{
		Subject:    subj,
		Name:       user.Name,
		Email:      user.Email,
		Gravatar:   user.Gravatar,
		Audience:   c.Name,
		Expiration: time.Now().Add(exp),
		IP:         user.IP,
		Location:   user.Location,
		IssuedAt:   time.Now(),
		Issuer:     issuerID,
	}

	jwt, err := claimid.ToJWT(privateKey)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	result := models.TokenResponse{
		Expiry: exp,
		Type:   "Bearer",
		Value:  jwt,
	}

	return http.StatusOK, result
}

func (t Tokening) GrantAuthcode() (int, interface{}) {
	return http.StatusNotImplemented, nil
}
