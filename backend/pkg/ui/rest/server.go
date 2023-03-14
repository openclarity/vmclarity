package rest

import "github.com/openclarity/vmclarity/shared/pkg/backendclient"

type ServerImpl struct {
	BackendClient *backendclient.BackendClient
}
