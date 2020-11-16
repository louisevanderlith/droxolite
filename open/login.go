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

func NewUILock(p *oidc.Provider, cfg *oauth2.Config) uiprotector {
	return uiprotector{authConfig: cfg, provider: p}
}

func (p uiprotector) Login(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	http.Redirect(w, r, p.authConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (p uiprotector) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	state, err := r.Cookie("oauthstate")

	if err != nil {
		http.Error(w, "state not found", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := p.authConfig.Exchange(ctx, r.URL.Query().Get("code"))
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

	location, err := r.Cookie("location")

	if err != nil {
		http.Error(w, "last location not found", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, location.Value, http.StatusFound)
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

func (p uiprotector) Middleware(next http.Handler) http.Handler {

	oidcConfig := &oidc.Config{
		ClientID: p.authConfig.ClientID,
	}
	v := p.provider.Verifier(oidcConfig)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawIDToken, err := r.Cookie("idtoken")

		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "location",
				Value:    r.RequestURI,
				Domain:   r.Host,
				Expires:  time.Now().Add(5 * time.Minute),
				Secure:   false,
				HttpOnly: true,
			})

			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		idToken, err := v.Verify(r.Context(), rawIDToken.Value)
		if err != nil {
			log.Println("Verify Error", err)
			http.SetCookie(w, &http.Cookie{
				Name:     "location",
				Value:    r.RequestURI,
				Domain:   r.Host,
				Expires:  time.Now().Add(5 * time.Minute),
				Secure:   false,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		rawaccToken, _ := r.Cookie("acctoken")
		xidn := context.WithValue(r.Context(), "Token", rawaccToken.Value)

		idn := context.WithValue(xidn, "IDToken", idToken)

		next.ServeHTTP(w, r.WithContext(idn))
	})
}
