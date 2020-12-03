package mix

import (
	"github.com/coreos/go-oidc"
	"net/http"
)

//ModFunc can be used to add additional values to the view bag before rendering
type ModFunc func(b Bag, r *http.Request)

func EndpointMod(endpoints map[string]string) ModFunc {
	return func(b Bag, r *http.Request) {
		b.SetValue("Endpoints", endpoints)
	}
}

func IdentityMod(clientId string) ModFunc {
	return func(b Bag, r *http.Request) {
		b.SetValue("ClientID", clientId)

		tkn := r.Context().Value("Token")

		if tkn == nil {
			return
		}

		b.SetValue("Token", tkn)

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

		b.SetValue("User", claims)
	}
}
