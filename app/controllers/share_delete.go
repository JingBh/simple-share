package controllers

import (
	"github.com/jingbh/simple-share/app/context"
	"github.com/jingbh/simple-share/internal/oss"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ShareDelete(c echo.Context) error {
	cc := c.(context.CustomContext)

	err := oss.DeleteShare(c.Request().Context(), cc.Share.Name)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
