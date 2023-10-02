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

package authstore

import "github.com/spf13/viper"

const (
	zitadelIssuerEnvVar        = "ZITADEL_ISSUER"
	zitadelAPIEnvVar           = "ZITADEL_API"
	zitadelInsecureEnvVar      = "ZITADEL_INSECURE"
	zitadelProjectIDEnvVar     = "ZITADEL_PROJECT_ID"
	zitadelOrgIDEnvVar         = "ZITADEL_ORG_ID"
	zitadelAuthKeyPathIDEnvVar = "ZITADEL_AUTH_KEY_PATH"
)

type Config struct {
	Issuer      string `json:"issuer,omitempty"`
	API         string `json:"api,omitempty"`
	Insecure    bool   `json:"insecure,omitempty"`
	ProjectID   string `json:"project-id,omitempty"`
	OrgID       string `json:"org-id,omitempty"`
	AuthKeyPath string `json:"auth-key-path,omitempty"`
}

func LoadConfig() Config {
	v := viper.New()
	v.AutomaticEnv()

	return Config{
		Issuer:      v.GetString(zitadelIssuerEnvVar),
		API:         v.GetString(zitadelAPIEnvVar),
		Insecure:    v.GetBool(zitadelInsecureEnvVar),
		ProjectID:   v.GetString(zitadelProjectIDEnvVar),
		OrgID:       v.GetString(zitadelOrgIDEnvVar),
		AuthKeyPath: v.GetString(zitadelAuthKeyPathIDEnvVar),
	}
}
