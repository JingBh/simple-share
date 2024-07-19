package utils

import (
	"github.com/spf13/viper"
	"strings"
)

func Url(path string) string {
	base := viper.GetString("baseurl")
	if strings.HasSuffix(base, "/") {
		base = base[:len(base)-1]
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}
