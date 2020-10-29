package open

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/coreos/go-oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type uiprotector struct {
	authConfig *oauth2.Config
}

func NewUILock(cfg *oauth2.Config) uiprotector {
	return uiprotector{authConfig: cfg}
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

	acccookie := http.Cookie{Name: "acctoken", Value: oauth2Token.AccessToken, Expires: oauth2Token.Expiry}
	http.SetCookie(w, &acccookie)

	idcookie := http.Cookie{Name: "idtoken", Value: rawIDToken, Expires: oauth2Token.Expiry}
	http.SetCookie(w, &idcookie)

	state.Expires = time.Now().Add(time.Hour * -24)
	state.Value = ""
	http.SetCookie(w, state)

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

func LoginMiddleware(verifier *oidc.IDTokenVerifier, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawIDToken, err := r.Cookie("idtoken")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		idToken, err := verifier.Verify(r.Context(), rawIDToken.Value)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rawaccToken, _ := r.Cookie("acctoken")
		xidn := context.WithValue(r.Context(), "Token", rawaccToken.Value)

		idn := context.WithValue(xidn, "IDToken", idToken)

		next.ServeHTTP(w, r.WithContext(idn))
	}
}
