// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authn

import "github.com/spf13/viper"

const (
	oidcIssuerEnvVar        = "AUTH_OIDC_ISSUER"
	oidcClientIDEnvVar      = "AUTH_OIDC_CLIENT_ID"
	oidcClientSecretEnvVar  = "AUTH_OIDC_CLIENT_SECRET" // #nosec G101
	oidcTokenURLEnvVar      = "AUTH_OIDC_TOKEN_URL"     // #nosec G101
	oidcIntrospectURLEnvVar = "AUTH_OIDC_INTROSPECT_URL"
	oidcUseZitadel          = "AUTH_USE_ZITADEL"
)

type Config struct {
	Issuer        string `json:"issuer,omitempty"`
	ClientID      string `json:"client-id,omitempty"`
	ClientSecret  string `json:"client-secret,omitempty"`
	TokenURL      string `json:"token-url,omitempty"`
	IntrospectURL string `json:"introspect-url,omitempty"`
	UseZitadel    bool   `json:"use-zitadel,omitempty"`
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
		UseZitadel:    v.GetBool(oidcUseZitadel),
	}
}
