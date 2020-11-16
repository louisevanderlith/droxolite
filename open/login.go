package open

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/coreos/go-oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

type uiprotector struct {
	provider   *oidc.Provider
	authConfig *oauth2.Config
}

func NewUILock(provider *oidc.Provider, cfg *oauth2.Config) uiprotector {
	return uiprotector{authConfig: cfg, provider: provider}
}

func (p uiprotector) Login(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	http.Redirect(w, r, p.authConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (p uiprotector) Callbackware(next http.Handler) http.Handler {
	v := p.provider.Verifier(&oidc.Config{
		ClientID: p.authConfig.ClientID,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		state, err := r.Cookie("oauthstate")

		if err != nil {
			http.Error(w, "state not found", http.StatusInternalServerError)
			return
		}

		if r.URL.Query().Get("state") != state.Value {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		state.Expires = time.Now().Add(time.Hour * -24)
		state.Value = ""
		state.HttpOnly = true
		http.SetCookie(w, state)

		oauth2Token, err := p.authConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)

		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}

		idToken, err := v.Verify(r.Context(), rawIDToken)

		if err != nil {
			log.Println("Verify Error", err)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		claims := make(map[string]interface{})
		err = idToken.Claims(&claims)

		if err != nil {
			log.Println("Claims Bind Error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		xidn := context.WithValue(r.Context(), "Token", oauth2Token)
		idn := context.WithValue(xidn, "Claims", claims)

		next.ServeHTTP(w, r.WithContext(idn))
	})
}

func (p uiprotector) Callback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("oauthstate")

	if err != nil {
		http.Error(w, "state not found", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := p.authConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)

	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	acccookie := http.Cookie{Name: "acctoken", Value: oauth2Token.AccessToken, Expires: oauth2Token.Expiry, HttpOnly: true}
	http.SetCookie(w, &acccookie)

	idcookie := http.Cookie{Name: "idtoken", Value: rawIDToken, Expires: oauth2Token.Expiry, HttpOnly: true}
	http.SetCookie(w, &idcookie)

	state.Expires = time.Now().Add(time.Hour * -24)
	state.Value = ""
	state.HttpOnly = true
	http.SetCookie(w, state)

	//TODO: Callback
	http.Redirect(w, r, "/", http.StatusFound)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
