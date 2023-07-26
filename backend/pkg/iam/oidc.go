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

package iam

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/backend/pkg/config"
	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"golang.org/x/oauth2"
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
//
// TODO: Enable support for creating ResourceServer for file data from string.
// TODO: This can be achieved using functional options, e.g. WithKey().
// TODO: Test against different OIDCs to check if this works. Tested against: Zitadel.
func NewOIDCProvider(config config.OIDC) (Provider, error) {
	// Add custom OIDC options
	var options []rs.Option
	if config.TokenURL != "" && config.IntrospectURL != "" {
		options = append(options, rs.WithStaticEndpoints("tokenUrl", "introspectUrl"))
	} else if config.TokenURL != "" || config.IntrospectURL != "" {
		return nil, fmt.Errorf("requires both OIDC token and introspect endpoints")
	}

	// Create resource server which provides introspection functionality
	var resourceServer rs.ResourceServer
	var err error
	if config.ClientKeyPath != "" {
		resourceServer, err = rs.NewResourceServerFromKeyFile(config.Issuer, config.ClientKeyPath, options...)
	} else {
		resourceServer, err = rs.NewResourceServerClientCredentials(config.Issuer, config.ClientID, config.ClientSecret, options...)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC resource server: %w", err)
	}

	// Return OIDC Provider
	return &oidcProvider{
		resourceServer: resourceServer,
		rolesClaim:     config.GetRolesClaim(),
	}, nil
}

func (provider *oidcProvider) Authenticate(ctx context.Context, request *http.Request) (*User, error) {
	// TODO: Explore caching options to reduce checks against identity server

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
	token, err := rs.Introspect(ctx, provider.resourceServer, authParts[1])
	if err != nil || !token.IsActive() {
		return nil, fmt.Errorf("authorization token is invalid")
	}

	// Get user roles from token role claim
	userRoles := make(map[string]bool)
	if tokenRolesClaim := token.GetClaim(provider.rolesClaim); tokenRolesClaim != nil {
		if tokenRolesMap, ok := tokenRolesClaim.(map[string]interface{}); ok {
			for tokenRole := range tokenRolesMap {
				userRoles[tokenRole] = true
			}
		} else if tokenRolesSlice, ok := tokenRolesClaim.([]string); ok {
			for _, tokenRole := range tokenRolesSlice {
				userRoles[tokenRole] = true
			}
		}
	}

	// Return user
	return &User{
		ID:    token.GetSubject(),
		Roles: userRoles,
	}, nil
}

// NewOIDCAuthInjectorFromKey creates a client injector which creates OAuth2 token
// source from key file to generate tokens that are injected in requests.
//
// TODO: Add support for injectors from file data string and Personal Access Tokens
func NewOIDCAuthInjectorFromKey(issuer, keyPath string, scopes []string) (client.RequestEditorFn, error) {
	// Get token source and token
	tokenSource, err := middleware.JWTProfileFromPath(keyPath)(issuer, append(scopes, oidc.ScopeOpenID))
	if err != nil {
		return nil, fmt.Errorf("unable to create OIDC token source: %w", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch OIDC token: %w", err)
	}

	// Use reusable token source to prevent request spikes
	tokenSource = oauth2.ReuseTokenSource(token, tokenSource)

	// Return client token injector function
	return func(_ context.Context, request *http.Request) error {
		token, err := tokenSource.Token()
		if err != nil {
			return err
		}

		token.SetAuthHeader(request)
		return nil
	}, nil
}
