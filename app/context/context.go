package context

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jingbh/simple-share/internal/models"
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	Token    *oidc.IDToken
	Username string
	Share    *models.Share
	echo.Context
}

func ExtractContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := ExtractToken(c.Request())

		cc := CustomContext{
			Token:    token,
			Username: ExtractUsername(token),
			Share:    ExtractShare(c),
			Context:  c,
		}
		return next(cc)
	}
}
