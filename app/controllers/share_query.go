package controllers

import (
	"github.com/jingbh/simple-share/internal/models"
	"github.com/jingbh/simple-share/internal/oss"
	"github.com/labstack/echo/v4"
)

type ShareListResponse struct {
	Shares     []*models.Share `json:"data"`
	NextCursor string          `json:"cursor"`
}

func ShareList(c echo.Context) error {
	cursor := c.QueryParam("cursor")
	res, nextCursor, err := oss.ListShares(c.Request().Context(), cursor)
	if err != nil {
		return err
	}

	return c.JSON(200, ShareListResponse{
		Shares:     res,
		NextCursor: nextCursor,
	})
}

func ShareGet(c echo.Context) error {
	name := c.Param("id")
	share, err := oss.GetShare(c.Request().Context(), name)
	if err != nil {
		return err
	}

	return c.JSON(200, share)
}
