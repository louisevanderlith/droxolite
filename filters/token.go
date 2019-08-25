package filters

import (
	"errors"
	"log"
	"strings"

	"github.com/louisevanderlith/droxolite/bodies"
	"github.com/louisevanderlith/droxolite/context"
	"github.com/louisevanderlith/droxolite/roletype"
)

func TokenCheck(ctx context.Contexer, requiredRole roletype.Enum, publicKeyPath, serviceName string) (bool, *bodies.Cookies) {
	path := ctx.RequestURI()
	//action := ctrl.ctx.Method()

	if strings.HasPrefix(path, "/favicon") {
		return true, nil
	}

	//requiredRole, err := m.GetRequiredRole(path, action)

	//if err != nil {
	//Missing Mapping, the user doesn't have access to the application
	//	ctx.RenderMethodResult(RenderUnauthorized(err))
	//	return
	//}

	if requiredRole == roletype.Unknown {
		//When we don't need credentials, we don't load it
		return true, nil
	}

	token, err := getAuthorizationToken(ctx)

	if err != nil {
		log.Println(err)
		//ctx.RenderMethodResult(RenderUnauthorized(err))
		return false, nil
	}

	avoc, err := bodies.GetAvoCookie(token, publicKeyPath)

	if err != nil {
		log.Println(err)
		//ctx.RenderMethodResult(RenderUnauthorized(err))
		return false, nil
	}

	allowed, err := bodies.IsAllowed(serviceName, avoc.UserRoles, requiredRole)

	if err != nil {
		log.Println(err)
		//ctx.RenderMethodResult(RenderUnauthorized(err))
		return false, nil
	}

	return allowed, avoc
}

//Returns the [TOKEN] in 'Bearer [TOKEN]'
func getAuthorizationToken(ctx context.Contexer) (string, error) {
	authHead, err := ctx.GetHeader("Authorization")

	if err != nil {
		return "", err
	}

	parts := strings.Split(authHead, " ")
	tokenType := parts[0]
	if strings.Trim(tokenType, " ") != "Bearer" {
		return "", errors.New("Bearer Authentication only")
	}

	return parts[1], nil
}
