package router

import (
	"authenticator-backend/config"
	"authenticator-backend/presentation/http/echo/handler"
	custom_middleware "authenticator-backend/presentation/http/echo/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

// SetRouter
// Summary: This is the function which sets the router for the echo.
// input: e(*echo.Echo): echo
// input: h(handler.AppHandler): handler
// input: config(*config.Config): config
// input: conn(*gorm.DB): gorm
// input: authMiddleware(custom_middleware.AuthMiddleware): auth middleware
func SetRouter(e *echo.Echo, h handler.AppHandler, config *config.Config, conn *gorm.DB, authMiddleware custom_middleware.AuthMiddleware) {
	if config.Env == "local" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,
		}))
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
		}))
	}

	e.Use(middleware.BodyLimit("25M"))

	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	e.GET("/api/v1/authInfo/health", func(c echo.Context) error { return h.HealthCheck(c) })

	authGroup := e.Group("")
	authGroup.Use(custom_middleware.APIKeyValidator(conn))
	if config.EnableIpRestriction {
		authGroup.Use(custom_middleware.IPForAPIKeyValidator(conn))
	}
	authGroup.PUT("/dataReset", func(c echo.Context) error { return h.Reset(c) }, authMiddleware.AuthJWT())

	auth := authGroup.Group("/auth")
	auth.Use(custom_middleware.AuthDump())
	auth.POST("/login", func(c echo.Context) error { return h.Login(c) })
	auth.POST("/refresh", func(c echo.Context) error { return h.Refresh(c) })
	auth.POST("/change", func(c echo.Context) error { return h.ChangePassword(c) }, authMiddleware.AuthJWT())

	systemAuth := authGroup.Group("/api/v1/systemAuth")
	systemAuth.Use(custom_middleware.SystemAPIKeyValidator(conn))
	systemAuth.Use(custom_middleware.AuthDump())
	systemAuth.POST("/token", func(c echo.Context) error { return h.TokenIntrospection(c) })
	systemAuth.POST("/apiKey", func(c echo.Context) error { return h.ApiKey(c) })

	authInfo := authGroup.Group("/api/v1/authInfo")
	authInfo.Use(authMiddleware.AuthJWT())
	authInfo.GET("", func(c echo.Context) error { return h.GetAuthInfo(c) })
	authInfo.PUT("", func(c echo.Context) error { return h.PutAuthInfo(c) })
}
