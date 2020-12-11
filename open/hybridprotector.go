package open

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/louisevanderlith/droxolite/mix"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

func (p hybridprotector) Refresh(w http.ResponseWriter, r *http.Request) {
	jtoken, _ := r.Cookie("acctoken")

	if jtoken == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	tkn64, err := base64.URLEncoding.DecodeString(jtoken.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tknVal := oauth2.Token{}
	err = json.Unmarshal(tkn64, &tknVal)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tknVal.Valid() {
		mix.Write(w, mix.JSON(jtoken.Value))
		return
	}

	params := "grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s"
	payload := strings.NewReader(fmt.Sprintf(params, p.usrConfig.ClientID, p.usrConfig.ClientSecret, tknVal.RefreshToken))

	req, _ := http.NewRequest("POST", p.provider.Endpoint().TokenURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println("ReadAll Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = mix.Write(w, mix.JSON(body))

	if err != nil {
		log.Println("Serve Error", err)
	}
}

func (p hybridprotector) Lock(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idtkn := r.Context().Value("IDToken")

		if idtkn == nil {
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

	tkn64 := base64.URLEncoding.EncodeToString(jtoken)
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
		setLastLocationCookie(w, r.URL.EscapedPath())

		jtoken, _ := r.Cookie("acctoken")

		if jtoken == nil {
			tkn, err := p.clntConfig.Token(r.Context())

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			acc := context.WithValue(r.Context(), "Token", *tkn)
			next.ServeHTTP(w, r.WithContext(acc))
			return
		}

		tkn64, err := base64.URLEncoding.DecodeString(jtoken.Value)

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

		idn := context.WithValue(xidn, "IDToken", idToken)
		//TODO: Replace IDToken with Claims (User)

		next.ServeHTTP(w, r.WithContext(idn))
	})
}
