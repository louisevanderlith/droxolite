package filters

import (
	"log"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/roletype"
)

//TokenCookieCheck is used to filter incoming UI Requests
func TokenCookieCheck(ctx context.Contexer, requiredRole roletype.Enum, publicKeyPath, serviceName string) (bool, *bodies.Cookies) {
	path := ctx.RequestURI()

	if strings.HasPrefix(path, "/static") || strings.HasPrefix(path, "/favicon") {
		return true, nil
	}

	token := ctx.FindQueryParam("access_token")

	if token == "" {
		cookie, err := ctx.GetCookie("avosession")

		if err != nil {
			log.Println(err)
		}

		token = cookie.Value
	}

	avoc, err := bodies.GetAvoCookie(token, publicKeyPath)

	if err != nil {
		log.Println(err)
	}

	allowed, err := bodies.IsAllowed(serviceName, avoc.UserRoles, requiredRole)

	if err != nil || !allowed {
		log.Println(err)
		return false, nil
	}

	return true, avoc
}
