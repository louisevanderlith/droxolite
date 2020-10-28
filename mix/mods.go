package mix

import (
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
		tkn := r.Context().Value("Token")

		f.SetValue("ClientID", clientId)
		f.SetValue("Token", tkn)
	}
}
