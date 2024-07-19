package app

import (
	"github.com/jingbh/simple-share/app/context"
	"github.com/jingbh/simple-share/app/controllers"
	"github.com/jingbh/simple-share/app/middlewares"
	"github.com/jingbh/simple-share/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"net/url"
)

func RegisterRoutes(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(context.ExtractContext)

	g := e.Group("auth/")
	g.Use(middlewares.CheckOIDCEnabled)
	g.Use(middlewares.DisableCache)
	g.GET("login", controllers.AuthRedirectLogin)
	g.GET("callback", controllers.AuthCallback)
	g.GET("userinfo", controllers.AuthGetUserinfo)

	g = e.Group("api/")
	g.GET("shares/:id", controllers.ShareGet, middlewares.ShareAuthenticated)
	// g.GET("shares/:id/content", controllers.ShareGetFile, middlewares.ShareAuthenticated)
	// g.GET("shares/:id/files/:fileId", controllers.ShareGetFile, middlewares.ShareAuthenticated)
	g.GET("shares", controllers.ShareList, middlewares.Authenticated)
	g.POST("shares", controllers.ShareCreate, middlewares.Authenticated)
	g.POST("upload", controllers.UploadStart, middlewares.Authenticated)
	g.POST("upload/:id/:part", controllers.UploadPart, middlewares.Authenticated)
	g.POST("upload/:id/complete", controllers.UploadComplete, middlewares.Authenticated)

	if viper.GetBool("embed.disable") {
		frontendUrl, err := url.Parse("http://localhost:5173")
		if err != nil {
			panic(err)
		}
		e.Any("/*", nil, middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{
			URL: frontendUrl,
		}})))
	} else {
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: web.HttpFs(),
		}))
	}
}
