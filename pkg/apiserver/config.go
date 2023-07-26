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

package apiserver

import (
	"encoding/json"
	"fmt"
	"strings"

	databaseTypes "github.com/openclarity/vmclarity/pkg/apiserver/database/types"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	OIDCIssuerEnvVar        = "OIDC_ISSUER"
	OIDCClientIDEnvVar      = "OIDC_CLIENT_ID"
	OIDCClientSecretEnvVar  = "OIDC_CLIENT_SECRET"
	OIDCClientKeyPathEnvVar = "OIDC_CLIENT_KEY_PATH"
	OIDCScopesEnvVar        = "OIDC_SCOPES"
	OIDCRolesClaimEnvVar    = "OIDC_ROLES_CLAIM"
	OIDCRolesClaimDefault   = "roles" // default role claim
	OIDCTokenURLEnvVar      = "OIDC_TOKEN_URL"
	OIDCIntrospectURLEnvVar = "OIDC_INTROSPECT_URL"

	BackendRestHost       = "BACKEND_REST_HOST"
	BackendRestDisableTLS = "BACKEND_REST_DISABLE_TLS" // nolint:gosec
	BackendRestPort       = "BACKEND_REST_PORT"
	HealthCheckAddress    = "HEALTH_CHECK_ADDRESS"

	DisableOrchestrator       = "DISABLE_ORCHESTRATOR"
	OrchestratorKeyPathEnvVar = "ORCHESTRATOR_KEY_PATH"

	UISitePath = "UI_SITE_PATH" // TODO: UI site should be moved out of the backend to nginx

	DBNameEnvVar     = "DB_NAME"
	DBUserEnvVar     = "DB_USER"
	DBPasswordEnvVar = "DB_PASS"
	DBHostEnvVar     = "DB_HOST"
	DBPortEnvVar     = "DB_PORT_NUMBER"
	DatabaseDriver   = "DATABASE_DRIVER"
	EnableDBInfoLogs = "ENABLE_DB_INFO_LOGS"

	OIDCIssuerEnvVar       = "OIDC_ISSUER"
	OIDCClientIDEnvVar     = "OIDC_CLIENT_ID"
	OIDCClientSecretEnvVar = "OIDC_CLIENT_SECRET"
	OIDCAppFilePathEnvVar  = "OIDC_APP_FILE_PATH"
	OIDCScopesEnvVar       = "OIDC_SCOPES"
	OIDCRolesClaimEnvVar   = "OIDC_ROLES_CLAIM"
	OIDCRolesClaimDefault  = "roles"

	OrchestratorKeyPathEnvVar = "ORCHESTRATOR_KEY_PATH"

	LocalDBPath = "LOCAL_DB_PATH"

	FakeDataEnvVar      = "FAKE_DATA"
	DisableOrchestrator = "DISABLE_ORCHESTRATOR"

	LogLevel = "LOG_LEVEL"
)

type Config struct {
	// Embed auth config
	OIDC

	// Backend
	BackendRestHost    string `json:"backend-rest-host,omitempty"`
	BackendRestPort    int    `json:"backend-rest-port,omitempty"`
	HealthCheckAddress string `json:"health-check-address,omitempty"`

	DisableOrchestrator bool   `json:"disable_orchestrator"`
	OrchestratorKeyPath string `json:"orchestrator-key-path"`

	// UI
	UISitePath string `json:"ui_site_path"`

	// Database config
	DatabaseDriver   string `json:"database-driver,omitempty"`
	DBName           string `json:"db-name,omitempty"`
	DBUser           string `json:"db-user,omitempty"`
	DBPassword       string `json:"-"`
	DBHost           string `json:"db-host,omitempty"`
	DBPort           string `json:"db-port,omitempty"`
	EnableDBInfoLogs bool   `json:"enable-db-info-logs"`
	EnableFakeData   bool   `json:"enable-fake-data"`
	LocalDBPath      string `json:"local-db-path,omitempty"`

	LogLevel log.Level `json:"log-level,omitempty"`
}

type OIDC struct {
	Issuer        string `json:"oidc-issuer"`
	ClientID      string `json:"oidc-client-id"`
	ClientSecret  string `json:"oidc-client-secret"`
	ClientKeyPath string `json:"oidc-client-key-path"`
	Scopes        string `json:"oidc-scopes"`
	RolesClaim    string `json:"oidc-roles-claim"`
	TokenURL      string `json:"oidc-token-url"`
	IntrospectURL string `json:"oidc-introspect-url"`
}

func (oidc *OIDC) GetScopes() []string {
	return strings.Split(oidc.Scopes, ",")
}

func (oidc *OIDC) GetRolesClaim() string {
	if oidc.RolesClaim == "" {
		return OIDCRolesClaimDefault
	}
	return oidc.RolesClaim
}

func setConfigDefaults() {
	viper.SetDefault(HealthCheckAddress, ":8081")
	viper.SetDefault(BackendRestPort, "8888")
	viper.SetDefault(DatabaseDriver, databaseTypes.DBDriverTypeLocal)
	viper.SetDefault(DisableOrchestrator, "false")

	viper.AutomaticEnv()
}

func LoadConfig() (*Config, error) {
	setConfigDefaults()

	config := &Config{}

	// Auth
	config.OIDC.Issuer = viper.GetString(OIDCIssuerEnvVar)
	config.OIDC.ClientID = viper.GetString(OIDCClientIDEnvVar)
	config.OIDC.ClientSecret = viper.GetString(OIDCClientSecretEnvVar)
	config.OIDC.ClientKeyPath = viper.GetString(OIDCClientKeyPathEnvVar)
	config.OIDC.Scopes = viper.GetString(OIDCScopesEnvVar)
	config.OIDC.RolesClaim = viper.GetString(OIDCRolesClaimEnvVar)
	config.OIDC.TokenURL = viper.GetString(OIDCTokenURLEnvVar)
	config.OIDC.IntrospectURL = viper.GetString(OIDCIntrospectURLEnvVar)

	// Backend
	config.BackendRestHost = viper.GetString(BackendRestHost)
	config.BackendRestPort = viper.GetInt(BackendRestPort)
	config.HealthCheckAddress = viper.GetString(HealthCheckAddress)

	config.DisableOrchestrator = viper.GetBool(DisableOrchestrator)
	config.OrchestratorKeyPath = viper.GetString(OrchestratorKeyPathEnvVar)

	// UI
	config.UISitePath = viper.GetString(UISitePath)

	// Database
	config.DatabaseDriver = viper.GetString(DatabaseDriver)
	config.DBPassword = viper.GetString(DBPasswordEnvVar)
	config.DBUser = viper.GetString(DBUserEnvVar)
	config.DBHost = viper.GetString(DBHostEnvVar)
	config.DBPort = viper.GetString(DBPortEnvVar)
	config.DBName = viper.GetString(DBNameEnvVar)
	config.EnableDBInfoLogs = viper.GetBool(EnableDBInfoLogs)
	config.EnableFakeData = viper.GetBool(FakeDataEnvVar)
	config.LocalDBPath = viper.GetString(LocalDBPath)

	// Common
	logLevel, err := log.ParseLevel(viper.GetString(LogLevel))
	if err != nil {
		logLevel = log.WarnLevel
	}
	config.LogLevel = logLevel

	configB, err := json.Marshal(config)
	if err == nil {
		log.Infof("\n\nconfig=%s\n\n", configB)
	} else {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	return config, nil
}
