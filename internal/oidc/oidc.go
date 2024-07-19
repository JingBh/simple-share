package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jingbh/simple-share/internal/utils"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io"
	"log"
)

var Enabled = false
var Provider *oidc.Provider
var Config *oidc.Config
var OAuthConfig *oauth2.Config

func GenerateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func InitOIDC() {
	var err error

	issuer := viper.GetString("oidc.issuer")
	clientId := viper.GetString("oidc.client_id")
	clientSecret := viper.GetString("oidc.client_secret")
	if issuer == "" {
		log.Println("OIDC is not configured")
		return
	}

	Provider, err = oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		log.Println("Failed to configure OIDC: ", err)
		return
	}

	Config = &oidc.Config{
		ClientID: clientId,
	}

	OAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     Provider.Endpoint(),
		RedirectURL:  utils.Url("/auth/callback"),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	Enabled = true
}
