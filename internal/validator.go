package internal

import (
	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"net/http"
)

type customValidator struct{}

func (cv *customValidator) Validate(i interface{}) error {
	v := validate.New(i)
	if v.Validate() {
		return nil
	}
	return &echo.HTTPError{
		Code:     http.StatusUnprocessableEntity,
		Message:  v.Errors.All(),
		Internal: v.Errors,
	}
}

func init() {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})
}
