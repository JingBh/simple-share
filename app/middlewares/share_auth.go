package middlewares

import (
	"github.com/jingbh/simple-share/app/context"
	"github.com/jingbh/simple-share/internal/oss"
	"github.com/jingbh/simple-share/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ShareAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(context.CustomContext)

		name := c.Param("id")
		share, _ := oss.GetShareCached(c.Request().Context(), name)
		if share == nil {
			return next(c)
		}

		if cc.Token != nil && share.Creator != nil && share.Creator.Subject == cc.Token.Subject {
			// is owner, skip authentication
			return next(c)
		}

		if share.Password != "" {
			password := c.Request().Header.Get("X-Share-Password")
			if password == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "this share requires password to access")
			}
			err := utils.VerifyPassword(password, share.Password)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
			}
		}

		return next(c)
	}
}
