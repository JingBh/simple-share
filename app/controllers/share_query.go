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

type ShareGetFileTypeResponse struct {
	Id   string          `json:"id"`
	Type models.FileType `json:"type"`
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
	requestHeaders.Add("Content-Type", contentType)
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
	}(res.Body)

	if v := res.Headers.Get("Cache-Control"); v != "" {
		c.Response().Header().Add("Cache-Control", v)
	}
	if v := res.Headers.Get("Content-Disposition"); v != "" {
		c.Response().Header().Add("Content-Disposition", v)
	}
	if v := res.Headers.Get("Content-Length"); v != "" {
		c.Response().Header().Add("Content-Length", v)
	}
	if v := res.Headers.Get("Content-Type"); v != "" {
		contentType = v
	}

	if c.Request().Method == http.MethodHead {
		return c.NoContent(res.StatusCode)
	}
	return c.Stream(res.StatusCode, contentType, res.Body)
}

func ShareGetFileType(c echo.Context) error {
	cc := c.(context.CustomContext)
	fileId := cc.Param("file")

	res, err := oss.GetShareContentType(c.Request().Context(), cc.Share.Name, fileId)
	if err != nil {
		return c.JSON(http.StatusOK, ShareGetFileTypeResponse{
			Id:   fileId,
			Type: models.FileTypeUnknown,
		})
	}
	c.Response().Header().Add("Cache-Control", "private, max-age=86400")
	return c.JSON(http.StatusOK, ShareGetFileTypeResponse{
		Id:   fileId,
		Type: res,
	})
}

func ShareGetFilePreview(c echo.Context) error {
	cc := c.(context.CustomContext)
	fileId := cc.Param("file")

	filetype, _ := oss.GetShareContentType(c.Request().Context(), cc.Share.Name, fileId)
	switch filetype {
	case models.FileTypeText:
		res, _ := oss.GetShareContent(c.Request().Context(), oss.GetShareContentOptions{
			Name:    cc.Share.Name,
			FileId:  fileId,
			Headers: nil,
		})
		if res != nil {
			defer func(reader io.ReadCloser) {
				_ = reader.Close()
			}(res.Body)
			return c.Stream(http.StatusOK, echo.MIMETextPlain, res.Body)
		}
		break
	case models.FileTypeImage:
		headers := make(http.Header)
		headers.Add("X-OSS-Process", "image/auto-orient,1/resize,m_lfit,l_2048,s_1536,limit_1/quality,q_90")
		res, _ := oss.GetShareContent(c.Request().Context(), oss.GetShareContentOptions{
			Name:    cc.Share.Name,
			FileId:  fileId,
			Headers: headers,
		})
		if res != nil {
			defer func(reader io.ReadCloser) {
				_ = reader.Close()
			}(res.Body)
			c.Response().Header().Add("Cache-Control", "private, max-age=86400")
			return c.Stream(http.StatusOK, "image/webp", res.Body)
		}
		break
	default:
		break
	}

	return echo.NewHTTPError(http.StatusNotFound, "preview not available")
}
