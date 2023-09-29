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
	os.Setenv("AUTH_OIDC_ISSUER", "http://localhost:8080")
	os.Setenv("AUTH_OIDC_CLIENT_ID", "233887130835812355@vmclarity")
	os.Setenv("AUTH_OIDC_CLIENT_SECRET", "OBEPiLBfjFZmZ7OzUGeFuOTlazq1yUqKSvsebhrZ8uOUCzU2jkdzksFH7l7T6pSB")
	os.Setenv("AUTH_OIDC_TOKEN_URL", "")
	os.Setenv("AUTH_OIDC_INTROSPECT_URL", "")

	// RUN
	Run(ctx, &Config{
		IamEnabled:          true,
		BackendRestPort:     8888,
		HealthCheckAddress:  ":8081",
		DisableOrchestrator: false,
		EnableFakeData:      false,
		DatabaseDriver:      "LOCAL",
		LocalDBPath:         "vmclarity.db",
	})
}
