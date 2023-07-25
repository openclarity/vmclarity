package iam

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/backend/pkg/config"
	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/oidc/pkg/oidc"
	"net/http"
	"strings"
)

type oidcProvider struct {
	resourceServer rs.ResourceServer
	rolesClaim     string
}

// NewOIDCProvider creates a Provider which intercepts requests and checks for a
// correct Bearer token using OAuth2 introspection by sending the token to the
// introspection endpoint.
func NewOIDCProvider(config *config.Config) (Provider, error) {
	var resourceServer rs.ResourceServer
	var err error
	if config.OIDCAppFilePath != "" {
		resourceServer, err = rs.NewResourceServerFromKeyFile(config.OIDCIssuer, config.OIDCAppFilePath)
	} else {
		resourceServer, err = rs.NewResourceServerClientCredentials(config.OIDCIssuer, config.OIDCClientID, config.OIDCClientSecret)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC IAM provider: %w", err)
	}

	return &oidcProvider{
		resourceServer: resourceServer,
		rolesClaim:     config.GetOIDCRolesClaim(),
	}, nil
}

func (provider *oidcProvider) Authenticate(ctx context.Context, request *http.Request) (*User, error) {
	// Validate authorization header
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}
	authParts := strings.Split(authHeader, oidc.PrefixBearer)
	if len(authParts) != 2 {
		return nil, fmt.Errorf("authorization header is malformed")
	}

	// Verify token against introspection endpoint
	resp, err := rs.Introspect(ctx, provider.resourceServer, authParts[1])
	if err != nil || !resp.IsActive() {
		return nil, fmt.Errorf("authorization token is invalid")
	}

	// Get user roles from token role claim
	roles := make(map[string]bool)
	if tokenRolesClaim, ok := resp.GetClaim(provider.rolesClaim).(map[string]interface{}); ok {
		for tokenRole := range tokenRolesClaim {
			roles[tokenRole] = true
		}
	}

	// Return user
	return &User{
		ID:    resp.GetSubject(),
		Roles: roles,
	}, nil
}
