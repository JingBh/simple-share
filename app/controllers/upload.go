package controllers

import (
	"github.com/jingbh/simple-share/internal/oss"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UploadStartResponse struct {
	FileId   string `json:"id"`
	PartSize int64  `json:"partSize"`
}

func UploadStart(c echo.Context) error {
	fileId, err := oss.UploadInit(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &UploadStartResponse{
		FileId:   fileId,
		PartSize: oss.UploadPartSize,
	})
}

func UploadPart(c echo.Context) error {
	fileId := c.Param("id")
	partNumber, err := strconv.Atoi(c.Param("part"))
	if err != nil || partNumber < 1 || partNumber > 10000 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid part number")
	}
	if c.Request().Body == nil || c.Request().ContentLength == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "empty body")
	}
	err = oss.UploadPart(c.Request().Context(), fileId, partNumber, c.Request().Body, c.Request().ContentLength)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func UploadComplete(c echo.Context) error {
	fileId := c.Param("id")
	err := oss.UploadComplete(c.Request().Context(), fileId)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
