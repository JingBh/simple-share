package context

import (
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	_oidc "github.com/jingbh/simple-share/internal/oidc"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func ExtractToken(req *http.Request) *oidc.IDToken {
	if !_oidc.Enabled {
		return nil
	}

	var rawToken string

	if header := req.Header.Get("Authorization"); header != "" {
		rawToken = strings.TrimPrefix(header, "Bearer ")
	} else {
		cookie, _ := req.Cookie("_token")
		if cookie != nil {
			rawToken = cookie.Value
		}
	}

	verifier := _oidc.Provider.Verifier(_oidc.Config)
	token, _ := verifier.Verify(req.Context(), rawToken)

	return token
}

func ExtractUsername(token *oidc.IDToken) string {
	if token == nil {
		return ""
	}

	var claims map[string]json.RawMessage
	if err := token.Claims(&claims); err != nil {
		return ""
	}

	var value string
	err := json.Unmarshal(claims[viper.GetString("oidc.name_claim")], &value)
	if err != nil {
		return ""
	}

	return value
}
