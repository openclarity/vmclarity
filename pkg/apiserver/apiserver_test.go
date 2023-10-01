package apiserver

import (
	"context"
	"github.com/openclarity/vmclarity/pkg/shared/log"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestAPI(t *testing.T) {
	runAPI() // run API

	// do API request
}

func runAPI() {
	ctx := context.Background()
	logger := logrus.WithContext(ctx)
	ctx = log.SetLoggerForContext(ctx, logger)

	// SET IAM CONFIG
	// AUTHN
	os.Setenv("AUTH_OIDC_ISSUER", "http://localhost:8080")
	os.Setenv("AUTH_OIDC_CLIENT_ID", "233887130835812355@vmclarity")
	os.Setenv("AUTH_OIDC_CLIENT_SECRET", "OBEPiLBfjFZmZ7OzUGeFuOTlazq1yUqKSvsebhrZ8uOUCzU2jkdzksFH7l7T6pSB")
	os.Setenv("AUTH_OIDC_TOKEN_URL", "")
	os.Setenv("AUTH_OIDC_INTROSPECT_URL", "")
	os.Setenv("AUTH_USE_ZITADEL", "true")

	// AUTHZ
	os.Setenv("AUTHZ_LOCAL_RBAC_RULE_FILEPATH", "/Users/rpolic/Desktop/repos/vmclarity/pkg/apiserver/iam/authz/rbac_rule_policy.csv")

	// AUTHSTORE
	os.Setenv("ZITADEL_ISSUER", "http://localhost:8080")
	os.Setenv("ZITADEL_API", "localhost:8080")
	os.Setenv("ZITADEL_INSECURE", "true")
	os.Setenv("ZITADEL_PROJECT_ID", "233887130080772099")
	os.Setenv("ZITADEL_ORG_ID", "233887129829113859")
	os.Setenv("ZITADEL_AUTH_KEY_PATH", "/Users/rpolic/Desktop/repos/vmclarity/pkg/apiserver/iam/testdata/zitadel/bootstrap/secrets/zitadel-admin-sa.json")

	// RUN
	Run(ctx, &Config{
		BackendRestPort:     8888,
		HealthCheckAddress:  ":8081",
		DisableOrchestrator: false,
		EnableFakeData:      false,
		DatabaseDriver:      "LOCAL",
		LocalDBPath:         "vmclarity.db",
	})
}
