package middlewares

import (
	"github.com/jingbh/simple-share/app/context"
	_oidc "github.com/jingbh/simple-share/internal/oidc"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckOIDCEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !_oidc.Enabled {
			return echo.NewHTTPError(http.StatusNotFound, "authentication is not configured")
		}
		return next(c)
	}
}

func Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(context.CustomContext)
		if cc.Token == nil {
			return echo.NewHTTPError(http.StatusForbidden)
		}
		return next(c)
	}
}
