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
	IAMEnabledEnvVar = "IAM_ENABLED"

	AuthOIDCIssuerEnvVar        = "AUTH_OIDC_ISSUER"
	AuthOIDCClientIDEnvVar      = "AUTH_OIDC_CLIENT_ID"
	AuthOIDCClientSecretEnvVar  = "AUTH_OIDC_CLIENT_SECRET"
	AuthOIDCTokenURLEnvVar      = "AUTH_OIDC_TOKEN_URL"
	AuthOIDCIntrospectURLEnvVar = "AUTH_OIDC_INTROSPECT_URL"

	AuthRoleSyncJwtRoleClaimEnvVar = "AUTH_ROLE_SYNC_JWT_ROLE_CLAIM"

	AuthorizationRbacRuleFilePathEnvVar = "AUTHZ_RBAC_RULE_FILE_PATH"

	BackendRestHost       = "BACKEND_REST_HOST"
	BackendRestDisableTLS = "BACKEND_REST_DISABLE_TLS" // nolint:gosec
	BackendRestPort       = "BACKEND_REST_PORT"
	HealthCheckAddress    = "HEALTH_CHECK_ADDRESS"

	DisableOrchestrator               = "DISABLE_ORCHESTRATOR"
	OrchestratorAuthBearerTokenEnvVar = "ORCHESTRATOR_AUTH_BEARER_TOKEN"

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

type AuthenticationOIDC struct {
	Issuer        string `json:"issuer"`
	ClientID      string `json:"client-id"`
	ClientSecret  string `json:"client-secret"`
	TokenURL      string `json:"token-url"`
	IntrospectURL string `json:"introspect-url"`
}

type Authentication struct {
	OIDC AuthenticationOIDC `json:"oidc"` // iam.Provider - OpenID Connect
}

type AuthRoleSynchronization struct {
	JWTRoleClaim string `json:"jwt-role-claim"` // iam.RoleSyncer - JWT Role Claim
}

type Authorization struct {
	RBACRuleFilePath string `json:"rbac-local"` // iam.Authorizer - RBAC Local
}

type Config struct {
	// IAM config
	IamEnabled              bool                    `json:"iam-enabled"`
	Authentication          Authentication          `json:"authentication"`
	AuthRoleSynchronization AuthRoleSynchronization `json:"auth-role-synchronization"`
	Authorization           Authorization           `json:"authorization"`

	// Backend
	BackendRestHost    string `json:"backend-rest-host,omitempty"`
	BackendRestPort    int    `json:"backend-rest-port,omitempty"`
	HealthCheckAddress string `json:"health-check-address,omitempty"`

	DisableOrchestrator         bool   `json:"disable_orchestrator"`
	OrchestratorAuthBearerToken string `json:"orchestrator-auth-bearer-token"`

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

type AuthenticationOIDC struct {
	Issuer        string `json:"issuer"`
	ClientID      string `json:"client-id"`
	ClientSecret  string `json:"client-secret"`
	ClientKeyPath string `json:"client-key-path"`
	TokenURL      string `json:"token-url"`
	IntrospectURL string `json:"introspect-url"`
}

type AuthInjectionJwt struct {
	Issuer      string `json:"issuer"`
	KeyPath     string `json:"key-path"`
	ExtraScopes string `json:"extra-scopes"` // defines additional scopes to fetch defined as a comma-separated string
}

func (jwt *AuthInjectionJwt) GetExtraScopes() []string {
	return strings.Split(jwt.ExtraScopes, ",")
}

type AuthorizationRBACLocal struct {
	RuleFilePath string `json:"rule-file-path"`
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

	// IAM
	config.IamEnabled = viper.GetBool(IAMEnabledEnvVar)

	// Auth - OIDC
	oidc := &config.Authentication.OIDC
	oidc.Issuer = viper.GetString(AuthOIDCIssuerEnvVar)
	oidc.ClientID = viper.GetString(AuthOIDCClientIDEnvVar)
	oidc.ClientSecret = viper.GetString(AuthOIDCClientSecretEnvVar)
	oidc.TokenURL = viper.GetString(AuthOIDCTokenURLEnvVar)
	oidc.IntrospectURL = viper.GetString(AuthOIDCIntrospectURLEnvVar)

	// AuthRoleSynchronization - JWT Role Claim
	syncRole := &config.AuthRoleSynchronization
	syncRole.JWTRoleClaim = viper.GetString(AuthRoleSyncJwtRoleClaimEnvVar)

	// Authorization - RBAC Local
	authzRbacLocal := &config.Authorization
	authzRbacLocal.RBACRuleFilePath = viper.GetString(AuthorizationRbacRuleFilePathEnvVar)

	// Backend
	config.BackendRestHost = viper.GetString(BackendRestHost)
	config.BackendRestPort = viper.GetInt(BackendRestPort)
	config.HealthCheckAddress = viper.GetString(HealthCheckAddress)

	config.DisableOrchestrator = viper.GetBool(DisableOrchestrator)
	config.OrchestratorAuthBearerToken = viper.GetString(OrchestratorAuthBearerTokenEnvVar)

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
