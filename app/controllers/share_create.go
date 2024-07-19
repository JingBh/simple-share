package controllers

import (
	_context "context"
	"encoding/json"
	"github.com/jingbh/simple-share/app/context"
	"github.com/jingbh/simple-share/internal/models"
	"github.com/jingbh/simple-share/internal/oss"
	"github.com/jingbh/simple-share/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

type shareCreateRequest struct {
	Type             string `json:"type" validate:"required|in:file,text,url"`
	Name             string `json:"name" validate:"required_if:nameRandom,false|nameValid" message:"required_if:please choose a name for your share|nameValid:invalid share name"`
	NameRandom       bool   `json:"nameRandom"`
	NameRandomLength int    `json:"nameRandomLength" validate:"required_if:nameRandom,true|range:4,32"`
	Password         string `json:"password" validate:"max_len:72"`
	Expiry           int    `json:"expiry" validate:"in:0,1,3,7"`
	Text             string `json:"text" validate:"required_unless:type,file|textIsUrl" message:"textIsUrl:invalid URL"`
	Files            []struct {
		Id   string `json:"id" validate:"required"`
		Path string `json:"path" validate:"required"`
	} `json:"files" validate:"required_if:type,file|max_len:100" message:"required_if:please upload at least one file first|max_len:too many files"`
}

func (r shareCreateRequest) NameValid(val string) bool {
	if !r.NameRandom {
		return oss.CheckShareName(val)
	}
	return true
}

func (r shareCreateRequest) TextIsUrl(val string) bool {
	if r.Type == "url" {
		_, err := url.Parse(val)
		return err == nil
	}
	return true
}

func ShareCreate(c echo.Context) error {
	cc := c.(context.CustomContext)
	req := new(shareCreateRequest)
	err := cc.Bind(req)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Internal: err,
		}
	}
	err = cc.Validate(req)
	if err != nil {
		return err
	}

	if req.NameRandom {
		req.Name, err = oss.GenerateShareName(req.NameRandomLength)
		if err != nil {
			return err
		}
	}

	var creator *models.ShareCreator
	if cc.Token != nil {
		creator = &models.ShareCreator{
			Subject:  cc.Token.Subject,
			Username: cc.Username,
		}
	}

	if req.Type == "text" || req.Type == "url" {
		err = oss.CreateShare(cc.Request().Context(), oss.CreateShareOptions{
			Type:     req.Type,
			Text:     req.Text,
			Path:     req.Name,
			Password: req.Password,
			Expiry:   req.Expiry,
			Creator:  creator,
		})
	} else if len(req.Files) > 1 {
		// directory
		var treeJsonBytes []byte
		treeJsonBytes, err = json.Marshal(req.Files)
		if err != nil {
			return err
		}
		treeJson := string(treeJsonBytes)
		for _, file := range req.Files {
			err = oss.CreateShare(cc.Request().Context(), oss.CreateShareOptions{
				Type:   "file",
				Source: file.Id + ".bin",
				Name:   utils.ExtractFilename(file.Path),
				Path:   req.Name + ".d/" + file.Id + ".bin",
				Expiry: req.Expiry,
			})
			if err != nil {
				break
			}
		}
		if err != nil {
			// error while creating files, the directory should be deleted to prevent orphan files.
			// the background context is used, so even if the request is cancelled, the deletion will still proceed.
			_ = oss.DeleteShare(_context.Background(), req.Name)
			return err
		}
		err = oss.CreateShare(cc.Request().Context(), oss.CreateShareOptions{
			Type:     "directory",
			Text:     treeJson,
			Path:     req.Name,
			Password: req.Password,
			Expiry:   req.Expiry,
			Creator:  creator,
		})
	} else {
		// single file, copy that file to destination
		err = oss.CreateShare(cc.Request().Context(), oss.CreateShareOptions{
			Type:     "file",
			Source:   req.Files[0].Id + ".bin",
			Name:     utils.ExtractFilename(req.Files[0].Path),
			Path:     req.Name,
			Password: req.Password,
			Expiry:   req.Expiry,
			Creator:  creator,
		})
	}
	if err != nil {
		return err
	}

	return cc.JSON(http.StatusOK, map[string]interface{}{
		"name": req.Name,
	})
}
