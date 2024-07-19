package internal

import (
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

func InitConfig() {
	viper.AutomaticEnv()
	viper.BindEnv("baseurl", "APP_BASEURL")
	viper.BindEnv("debug", "APP_DEBUG")
	viper.BindEnv("serve.host", "HOST")
	viper.BindEnv("serve.port", "PORT")
	viper.BindEnv("oss.access_key_id", "ALIBABA_CLOUD_ACCESS_KEY_ID")
	viper.BindEnv("oss.access_key_secret", "ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("debug", false)
	viper.SetDefault("serve.port", 8080)
	viper.SetDefault("oidc.name_claim", "username")

	if viper.GetBool("debug") {
		viper.SetDefault("serve.host", "localhost")
		viper.SetDefault("embed.disable", true)
	} else {
		viper.SetDefault("serve.host", "")
		viper.SetDefault("embed.disable", false)
	}

	{
		host := viper.GetString("serve.host")
		port := viper.GetInt("serve.port")
		if host == "0.0.0.0" || host == "::" {
			host = ""
		}
		if strings.Contains(host, ":") {
			host = "[" + host + "]"
		}
		addr := host + ":" + strconv.Itoa(port)
		viper.Set("serve.addr", addr)
		viper.SetDefault("baseurl", "http://"+addr)
	}
}
