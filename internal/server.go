package internal

import (
	"errors"
	"github.com/jingbh/simple-share/app"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func StartServer() {
	e := echo.New()
	e.Debug = viper.GetBool("debug")
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.Validator = &customValidator{}

	app.RegisterRoutes(e)

	if err := e.Start(viper.GetString("serve.addr")); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
