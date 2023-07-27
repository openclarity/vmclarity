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
	"github.com/openclarity/vmclarity/backend/pkg/config"
	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type oidcInjector struct {
	tokenSource oauth2.TokenSource
}

// NewOIDCInjector creates a client Injector which creates OAuth2 token source
// from key file to generate tokens that will be injected in requests.
//
// TODO: Enable support for creating token sources from file data string and Personal Access Tokens.
// TODO: This can be achieved using functional options, e.g. WithKey(string), WithKeyFile(filepath), WithAccessToken(string).
// TODO: Test against different OIDCs to check if this works. Tested against: Zitadel.
func NewOIDCInjector(issuer, keyPath string, extraScopes []string) (Injector, error) {
	// Get token source
	tokenSource, err := middleware.JWTProfileFromPath(keyPath)(issuer, append(extraScopes, oidc.ScopeOpenID))
	if err != nil {
		return nil, fmt.Errorf("unable to create OIDC token source: %w", err)
	}

	// Return Injector with reusable token source to prevent request spikes
	return &oidcInjector{
		tokenSource: oauth2.ReuseTokenSource(nil, tokenSource),
	}, nil
}

func (injector *oidcInjector) Inject(_ context.Context, request *http.Request) error {
	token, err := injector.tokenSource.Token()
	if err != nil {
		return err
	}

	token.SetAuthHeader(request)
	return nil
}

type oidcProvider struct {
	resourceServer rs.ResourceServer
	roleClaim      string
}

// NewOIDCProvider creates a Provider which intercepts requests and checks for a
// correct Bearer token using OAuth2 introspection by sending the token to the
// introspection endpoint.
//
// TODO: Enable support for creating resource server from file data string.
// TODO: This can be achieved using functional options, e.g. WithKey(string), WithKeyFile(filepath).
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
		roleClaim:      config.GetRoleClaim(),
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
	token, err := rs.Introspect(ctx, provider.resourceServer, authParts[1])
	if err != nil || !token.IsActive() {
		return nil, fmt.Errorf("authorization token is invalid")
	}

	// Get user roles from token role claim
	var userRoles []string
	switch tokenRoles := token.GetClaim(provider.roleClaim).(type) {
	case map[string]interface{}:
		index := 0
		userRoles = make([]string, len(tokenRoles))
		for roleClaim := range tokenRoles {
			userRoles[index] = roleClaim
			index++
		}
	case []string:
		userRoles = tokenRoles
	}

	// Return user
	return &User{
		ID:    token.GetSubject(),
		Roles: userRoles,
	}, nil
}
