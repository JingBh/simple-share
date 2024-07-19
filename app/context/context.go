package context

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	Token    *oidc.IDToken
	Username string
	echo.Context
}

func ExtractContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := ExtractToken(c.Request())

		cc := CustomContext{
			Token:    token,
			Username: ExtractUsername(token),
			Context:  c,
		}
		return next(cc)
	}
}
