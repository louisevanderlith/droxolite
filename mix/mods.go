package mix

import (
	"github.com/coreos/go-oidc"
	"net/http"
)

type ModFunc func(f MixerFactory, r *http.Request)

func EndpointMod(endpoints map[string]string) ModFunc {
	return func(f MixerFactory, r *http.Request) {
		f.SetValue("Endpoints", endpoints)
	}
}

func IdentityMod(clientId string) ModFunc {
	return func(f MixerFactory, r *http.Request) {
		f.SetValue("ClientID", clientId)

		tkn := r.Context().Value("Token")

		if tkn == nil {
			return
		}

		f.SetValue("Token", tkn)

		tknVal := r.Context().Value("IDToken")

		if tknVal == nil {
			return
		}

		idToken := tknVal.(*oidc.IDToken)
		claims := make(map[string]interface{})
		err := idToken.Claims(&claims)

		if err != nil {
			panic(err)
			return
		}

		f.SetValue("User", claims)
	}
}
