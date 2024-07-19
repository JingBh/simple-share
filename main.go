package main

import (
	"github.com/jingbh/simple-share/internal"
	"github.com/jingbh/simple-share/internal/oidc"
)

func main() {
	internal.InitConfig()
	go func() {
		oidc.InitOIDC()
	}()
	internal.StartServer()
}
