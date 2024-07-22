package context

import (
	"github.com/jingbh/simple-share/internal/models"
	"github.com/jingbh/simple-share/internal/oss"
	"github.com/labstack/echo/v4"
	"strings"
)

func ExtractShare(c echo.Context) *models.Share {
	if !strings.Contains(c.Path(), "shares") {
		return nil
	}

	name := c.Param("name")
	if name != "" {
		share, _ := oss.GetShareCached(c.Request().Context(), name)
		return share
	}
	return nil
}
