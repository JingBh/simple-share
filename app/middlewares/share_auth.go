package middlewares

import (
	"github.com/jingbh/simple-share/app/context"
	"github.com/jingbh/simple-share/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ShareAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(context.CustomContext)

		if cc.Share == nil {
			return echo.NewHTTPError(http.StatusNotFound, "share not found")
		}

		if cc.Token != nil && cc.Share.Creator != nil && cc.Share.Creator.Subject == cc.Token.Subject {
			// is owner, skip authentication
			return next(c)
		}

		if cc.Share.Password != "" {
			password := c.QueryParam("password")
			if password == "" {
				password = c.Request().Header.Get("X-Share-Password")
			}
			if password == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "this share requires password to access")
			}
			err := utils.VerifyPassword(password, cc.Share.Password)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
			}
		}

		return next(c)
	}
}

func ShareAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(context.CustomContext)

		if cc.Share == nil {
			return echo.NewHTTPError(http.StatusNotFound, "share not found")
		}

		if cc.Token != nil && cc.Share.Creator != nil && cc.Share.Creator.Subject == cc.Token.Subject {
			// is owner, grant access
			return next(c)
		}

		return echo.NewHTTPError(http.StatusForbidden, "you are not authorized to do this operation on this share")
	}
}
