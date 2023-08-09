package utils

import (
	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
	"os"
)

const APITokenEnvVar = "VMCLARITY_API_TOKEN" // #nosec G101

// NewBackendClient creates a backendclient.BackendClient and optionally adds bearer token.
func NewBackendClient(server string) (*backendclient.BackendClient, error) {
	return backendclient.Create(server, backendclient.WithTokenAuth(os.Getenv(APITokenEnvVar)))
}
