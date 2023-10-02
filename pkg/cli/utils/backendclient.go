package utils

import (
	"os"

	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
)

const APITokenEnvVar = "VMCLARITY_API_TOKEN" // #nosec G101

// NewBackendClient creates a backendclient.BackendClient and optionally adds bearer token.
func NewBackendClient(server string) (*backendclient.BackendClient, error) {
	return backendclient.Create(server, backendclient.WithTokenAuth(os.Getenv(APITokenEnvVar))) // nolint:wrapcheck
}
