package open

import (
	"encoding/base64"
	"encoding/json"
	"github.com/coreos/go-oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"net/http"
)

func NewHybridLock(p *oidc.Provider, clntCfg *clientcredentials.Config, usrConfig *oauth2.Config) hybridprotector {
	return hybridprotector{
		provider:   p,
		clntConfig: clntCfg,
		usrConfig:  usrConfig,
	}
}

type hybridprotector struct {
	provider   *oidc.Provider
	clntConfig *clientcredentials.Config
	usrConfig  *oauth2.Config
}

func (p hybridprotector) Lock(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idtkn := r.Context().Value("IDToken")

		if idtkn == nil {
			setLastLocationCookie(w, r.URL.EscapedPath())
			p.Login(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (p hybridprotector) Login(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	http.Redirect(w, r, p.usrConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (p hybridprotector) Callback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("oauthstate")

	if err != nil {
		http.Error(w, "state not found", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := p.usrConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)

	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	jtoken, err := json.Marshal(oauth2Token)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tkn64 := base64.StdEncoding.EncodeToString(jtoken)
	tokencookie := http.Cookie{
		Name:     "acctoken",
		Value:    tkn64,
		MaxAge:   0,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &tokencookie)

	idcookie := http.Cookie{
		Name:     "idtoken",
		Value:    rawIDToken,
		MaxAge:   0,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &idcookie)

	state.MaxAge = -1
	state.Value = ""
	http.SetCookie(w, state)

	RedirectToLastLocation(w, r)
}

func (p hybridprotector) Logout(w http.ResponseWriter, r *http.Request) {
	acc, err := r.Cookie("acctoken")

	if err != nil {
		http.Error(w, "acctoken not found", http.StatusInternalServerError)
		return
	}

	acc.MaxAge = -1
	acc.Value = ""
	http.SetCookie(w, acc)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (p hybridprotector) Protect(next http.Handler) http.Handler {
	oidcConfig := &oidc.Config{
		ClientID: p.usrConfig.ClientID,
	}

	v := p.provider.Verifier(oidcConfig)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jtoken, _ := r.Cookie("acctoken")

		if jtoken == nil {
			tkn, err := p.clntConfig.Token(r.Context())

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			acc := context.WithValue(r.Context(), "Token", *tkn)
			next.ServeHTTP(w, r.WithContext(acc))
		}

		tkn64, err := base64.StdEncoding.DecodeString(jtoken.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accToken := oauth2.Token{}
		err = json.Unmarshal(tkn64, &accToken)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		xidn := context.WithValue(r.Context(), "Token", accToken)

		rawIDToken, err := r.Cookie("idtoken")

		if err != nil {
			log.Println("Cookie Error", err)
			next.ServeHTTP(w, r.WithContext(xidn))
			return
		}

		idToken, err := v.Verify(r.Context(), rawIDToken.Value)
		if err != nil {
			log.Println("ID Verify Error", err)
			next.ServeHTTP(w, r.WithContext(xidn))
			return
		}

		err = idToken.VerifyAccessToken(accToken.AccessToken)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//xidn := context.WithValue(r.Context(), "Token", accToken)
		idn := context.WithValue(xidn, "IDToken", idToken)
		//TODO: Replace IDToken with Claims (User)

		next.ServeHTTP(w, r.WithContext(idn))
	})
}
