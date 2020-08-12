package drx

import (
	"github.com/louisevanderlith/kong/tokens"
	"net/http"
)

//GetToken will return 'token' assigned to Context
func GetToken(r *http.Request) string {
	v := r.Context().Value("token")

	return v.(string)
}

func GetIdentity(r *http.Request) tokens.Identity {
	val := r.Context().Value("claims")

	res, ok := val.(tokens.Identity)

	if !ok {
		return nil
	}

	return res
}

func GetUserIdentity(r *http.Request) tokens.UserIdentity {
	val := r.Context().Value("userclaims")

	res, ok := val.(tokens.UserIdentity)

	if !ok {
		return nil
	}

	return res
}
