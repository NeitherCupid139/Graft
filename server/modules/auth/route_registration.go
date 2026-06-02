package auth

import (
	"github.com/gin-gonic/gin"

	authopenapi "graft/server/internal/contract/openapi/auth"
	"graft/server/internal/httpx"
	"graft/server/internal/module"
	"graft/server/internal/moduleapi"
	authcontract "graft/server/modules/auth/contract"
)

type routeGuards struct {
	authenticated          gin.HandlerFunc
	requiredPasswordChange gin.HandlerFunc
	restrictedSession      gin.HandlerFunc
}

type authRouteRegistrar struct {
	ctx        *module.Context
	pluginName string
	authFlow   moduleapi.AuthFlowService
	cookies    CookieManager
	guards     routeGuards
}

func registerAuthRoutes(
	ctx *module.Context,
	pluginName string,
	authService moduleapi.AuthService,
	authFlow moduleapi.AuthFlowService,
) error {
	authGroup := ctx.Router.Group(authcontract.AuthGroup)
	guards := newRouteGuards(ctx, authService, authFlow, authGroup.BasePath())

	registrar := authRouteRegistrar{
		ctx:        ctx,
		pluginName: pluginName,
		authFlow:   authFlow,
		cookies:    NewCookieManager(ctx.Config.Auth),
		guards:     guards,
	}
	authGroup.Use(httpx.RequestIDMiddleware())
	registrar.registerLoginRoutes(authGroup)
	registrar.registerCurrentUserSessionRoutes(authGroup)
	registrar.registerBootstrapAndPasswordRoutes(authGroup)

	return nil
}

var _ authopenapi.ServerInterface = authGeneratedHandler{}
