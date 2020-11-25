package open

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

func NewGhostware(cfg *clientcredentials.Config) ghostprotector {
	return ghostprotector{cfg: cfg}
}

type ghostprotector struct {
	cfg *clientcredentials.Config
}

func (g ghostprotector) GhostMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tkn, err := g.cfg.Token(r.Context())

		if err != nil {
			panic(err)
		}

		acc := context.WithValue(r.Context(), "Token", *tkn)
		next.ServeHTTP(w, r.WithContext(acc))
	}
}
