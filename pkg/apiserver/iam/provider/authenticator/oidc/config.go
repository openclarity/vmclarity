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

package oidc

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	authOIDCIssuerEnvVar        = "AUTH_OIDC_ISSUER"
	authOIDCClientIDEnvVar      = "AUTH_OIDC_CLIENT_ID"
	authOIDCClientSecretEnvVar  = "AUTH_OIDC_CLIENT_SECRET" // #nosec G101
	authOIDCTokenURLEnvVar      = "AUTH_OIDC_TOKEN_URL"     // #nosec G101
	authOIDCIntrospectURLEnvVar = "AUTH_OIDC_INTROSPECT_URL"
)

type Config struct {
	Issuer        string `json:"issuer"`
	ClientID      string `json:"client-id"`
	ClientSecret  string `json:"client-secret"`
	TokenURL      string `json:"token-url"`
	IntrospectURL string `json:"introspect-url"`
}

func LoadConfig() (*Config, error) {
	return &Config{
		Issuer:        viper.GetString(authOIDCIssuerEnvVar),
		ClientID:      viper.GetString(authOIDCClientIDEnvVar),
		ClientSecret:  viper.GetString(authOIDCClientSecretEnvVar),
		TokenURL:      viper.GetString(authOIDCTokenURLEnvVar),
		IntrospectURL: viper.GetString(authOIDCIntrospectURLEnvVar),
	}, nil
}

func (c *Config) Validate() error {
	// Validate params
	if c.Issuer == "" {
		return fmt.Errorf("oidc: must specify issuer")
	}
	if c.ClientID == "" {
		return fmt.Errorf("oidc: must specify client id")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("oidc: must specify client secret")
	}

	// Validate endpoints, it must specify either none or both endpoints
	hasTokenUrl := c.TokenURL != ""
	hasIntrospectUrl := c.IntrospectURL != ""
	if hasTokenUrl != hasIntrospectUrl {
		return fmt.Errorf("oidc: must specify both token and introspect endpoints")
	}

	return nil
}
