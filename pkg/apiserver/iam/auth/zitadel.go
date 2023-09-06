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

package auth

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"
	"github.com/zitadel/oidc/pkg/oidc"
	"net/http"
	"strings"

	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
)

// New creates an authenticator which intercepts requests and checks for a
// correct Bearer token using OAuth2 introspection by sending the token to the
// introspection endpoint. On success, returns an iam.User with configured
// JwtClaims.
func New(config types.AuthConfig) (types.Authenticator, error) {
	// Add custom OIDC options
	var options []rs.Option
	if config.TokenURL != "" && config.IntrospectURL != "" {
		options = append(options, rs.WithStaticEndpoints(config.TokenURL, config.IntrospectURL))
	}

	// Create resource server which provides introspection functionality
	resourceServer, err := rs.NewResourceServerClientCredentials(config.Issuer, config.ClientID, config.ClientSecret, options...)
	if err != nil {
		return nil, fmt.Errorf("could not create resource server for Authenticator=OIDC: %w", err)
	}

	// Return OIDC Authenticator
	return &authService{
		resourceServer: resourceServer,
		mgmtClient:     nil,
	}, nil
}

type authService struct {
	projectID      string
	resourceServer rs.ResourceServer
	mgmtClient     *management.Client
}

func (auth *authService) RevokeAccess(ctx context.Context, user types.User) error {
	return nil
}

func (auth *authService) Authenticate(ctx context.Context, req *http.Request) (*types.User, error) {
	// Extract token
	token, err := extractToken(req)
	if err != nil {
		return nil, err
	}

	// Verify token against introspection endpoint
	jwtToken, err := rs.Introspect(ctx, auth.resourceServer, token)
	if err != nil {
		return nil, fmt.Errorf("token invalid: %w", err)
	}
	if !jwtToken.IsActive() {
		return nil, fmt.Errorf("token expired")
	}

	// Load roles for given user
	roles, _ := jwtToken.GetClaim("roles").([]string)

	// Return user
	return &types.User{
		ID:    jwtToken.GetSubject(),
		Roles: roles,
	}, nil
}

func extractToken(r *http.Request) (string, error) {
	auth := r.Header.Get("authorization")
	if auth == "" {
		// http.StatusUnauthorized
		return "", fmt.Errorf("auth header missing")
	}
	if !strings.HasPrefix(auth, oidc.PrefixBearer) {
		// http.StatusUnauthorized
		return "", fmt.Errorf("invalid auth header")
	}
	return strings.TrimPrefix(auth, oidc.PrefixBearer), nil
}
