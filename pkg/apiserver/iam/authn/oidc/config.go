package oidc

import "github.com/spf13/viper"

const (
	oidcIssuerEnvVar        = "AUTH_OIDC_ISSUER"
	oidcClientIDEnvVar      = "AUTH_OIDC_CLIENT_ID"
	oidcClientSecretEnvVar  = "AUTH_OIDC_CLIENT_SECRET" // #nosec G101
	oidcTokenURLEnvVar      = "AUTH_OIDC_TOKEN_URL"     // #nosec G101
	oidcIntrospectURLEnvVar = "AUTH_OIDC_INTROSPECT_URL"
)

type Config struct {
	Issuer        string `json:"issuer,omitempty"`
	ClientID      string `json:"client-id,omitempty"`
	ClientSecret  string `json:"client-secret,omitempty"`
	TokenURL      string `json:"token-url,omitempty"`
	IntrospectURL string `json:"introspect-url,omitempty"`
}

func LoadConfig() Config {
	v := viper.New()
	v.AutomaticEnv()

	return Config{
		Issuer:        v.GetString(oidcIssuerEnvVar),
		ClientID:      v.GetString(oidcClientIDEnvVar),
		ClientSecret:  v.GetString(oidcClientSecretEnvVar),
		TokenURL:      v.GetString(oidcTokenURLEnvVar),
		IntrospectURL: v.GetString(oidcIntrospectURLEnvVar),
	}
}
