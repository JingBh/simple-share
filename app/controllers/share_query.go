package controllers

import (
	"github.com/jingbh/simple-share/app/context"
	"github.com/jingbh/simple-share/internal/models"
	"github.com/jingbh/simple-share/internal/oss"
	"github.com/jingbh/simple-share/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"io"
	"net/http"
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
	cc := c.(context.CustomContext)
	return c.JSON(200, cc.Share)
}

func ShareShow(c echo.Context) error {
	name := c.Param("name")
	return c.Redirect(http.StatusFound, utils.Url("/#/shares/"+name))
}

func ShareGetFile(c echo.Context) error {
	cc := c.(context.CustomContext)
	fileId := cc.Param("file")

	contentType := "application/octet-stream"
	if cc.Share.Type == "directory" && fileId == "" {
		contentType = "application/json"
	} else if cc.Share.Type == "text" || cc.Share.Type == "url" {
		contentType = "text/plain"
	}

	if viper.GetBool("oss.download_direct") {
		url, err := oss.GetShareContentLink(c.Request().Context(), oss.GetShareContentLinkOptions{
			Name:        cc.Share.Name,
			FileId:      fileId,
			ContentType: contentType,
		})
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, url)
	}

	c.Response().Header().Add("Accept-Ranges", "bytes")
	requestHeaders := make(http.Header)
	requestHeaders.Add("Range", c.Request().Header.Get("Range"))
	res, err := oss.GetShareContent(c.Request().Context(), oss.GetShareContentOptions{
		Name:    cc.Share.Name,
		FileId:  fileId,
		Headers: requestHeaders,
	})
	if err != nil {
		return err
	}
	if res == nil {
		return echo.NewHTTPError(http.StatusNotFound, "file not found")
	}
	defer func(reader io.ReadCloser) {
		_ = reader.Close()
	}(res.Reader)

	if v := res.Headers.Get("Cache-Control"); v != "" {
		c.Response().Header().Add("Cache-Control", v)
	}
	if v := res.Headers.Get("Content-Length"); v != "" {
		c.Response().Header().Add("Content-Length", v)
	}
	if v := res.Headers.Get("Content-Disposition"); v != "" {
		c.Response().Header().Add("Content-Disposition", v)
	}

	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}
	return c.Stream(http.StatusOK, contentType, res.Reader)
}
