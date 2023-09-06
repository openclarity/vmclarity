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

package types

import (
	"context"
	"github.com/spf13/viper"
	"net/http"
)

const (
	authOIDCIssuerEnvVar        = "AUTH_OIDC_ISSUER"
	authOIDCClientIDEnvVar      = "AUTH_OIDC_CLIENT_ID"
	authOIDCClientSecretEnvVar  = "AUTH_OIDC_CLIENT_SECRET" // #nosec G101
	authOIDCTokenURLEnvVar      = "AUTH_OIDC_TOKEN_URL"     // #nosec G101
	authOIDCIntrospectURLEnvVar = "AUTH_OIDC_INTROSPECT_URL"
)

type AuthConfig struct {
	Issuer        string `json:"issuer,omitempty"`
	ClientID      string `json:"client-id,omitempty"`
	ClientSecret  string `json:"client-secret,omitempty"`
	TokenURL      string `json:"token-url,omitempty"`
	IntrospectURL string `json:"introspect-url,omitempty"`
}

func LoadConfig() AuthConfig {
	return AuthConfig{
		Issuer:        viper.GetString(authOIDCIssuerEnvVar),
		ClientID:      viper.GetString(authOIDCClientIDEnvVar),
		ClientSecret:  viper.GetString(authOIDCClientSecretEnvVar),
		TokenURL:      viper.GetString(authOIDCTokenURLEnvVar),
		IntrospectURL: viper.GetString(authOIDCIntrospectURLEnvVar),
	}
}

type IAMService interface {
	Authenticator() Authenticator
	Authorizer() Authorizer
}

//type ScansTable interface {
//	GetScans(params models.GetScansParams) (models.Scans, error)
//	GetScan(scanID models.ScanID, params models.GetScansScanIDParams) (models.Scan, error)
//
//	CreateScan(scan models.Scan) (models.Scan, error)
//	UpdateScan(scan models.Scan, params models.PatchScansScanIDParams) (models.Scan, error)
//	SaveScan(scan models.Scan, params models.PutScansScanIDParams) (models.Scan, error)
//
//	DeleteScan(scanID models.ScanID) error
//}

// User defines an authenticated user.
type User struct {
	ID    string   `json:"ID"`
	Roles []string `json:"roles"`
}

// Authenticator defines authentication service.
type Authenticator interface {
	// Authenticate authenticates a request payload against authentication service.
	Authenticate(ctx context.Context, req *http.Request) (*User, error)
}

// Authorizer defines authorization service.
type Authorizer interface {
	// CanPerform checks if User is allowed to perform an action on an asset based on
	// some predefined rules.
	CanPerform(ctx context.Context, user *User, asset, action string) (bool, error)
}
