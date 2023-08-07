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

package oidc

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"net/http"
	"strings"

	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/oidc/pkg/oidc"
)

// New creates an authenticator which intercepts requests and checks for a
// correct Bearer token using OAuth2 introspection by sending the token to the
// introspection endpoint. On success, returns an iam.User with configured
// JwtClaims.
//
// TODO: Test against different OIDCs to check if this works. Tested against: Zitadel.
func New() (iam.Authenticator, error) {
	// Load config
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to load config: %w", err)
	}
	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("oidc: failed to validate config: %w", err)
	}

	// Add custom OIDC options
	var options []rs.Option
	if config.TokenURL != "" && config.IntrospectURL != "" {
		options = append(options, rs.WithStaticEndpoints(config.TokenURL, config.IntrospectURL))
	}

	// Create resource server which provides introspection functionality
	resourceServer, err := rs.NewResourceServerClientCredentials(config.Issuer, config.ClientID, config.ClientSecret, options...)
	if err != nil {
		return nil, fmt.Errorf("oidc: could not create resource server: %w", err)
	}

	// Return OIDC Authenticator
	return &oidcAuth{
		resourceServer: resourceServer,
	}, nil
}

type oidcAuth struct {
	resourceServer rs.ResourceServer
}

func (auth *oidcAuth) Authenticate(ctx context.Context, request *http.Request) (*iam.User, error) {
	// Validate authorization header
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("oidc: authorization header is missing")
	}
	authParts := strings.Split(authHeader, oidc.PrefixBearer)
	if len(authParts) != 2 {
		return nil, fmt.Errorf("oidc: authorization header is malformed")
	}

	// Verify token against introspection endpoint
	token, err := rs.Introspect(ctx, auth.resourceServer, authParts[1])
	if err != nil || !token.IsActive() {
		return nil, fmt.Errorf("oidc: authorization token is invalid")
	}

	// Return user
	return &iam.User{
		ID:        token.GetSubject(),
		JwtClaims: token.GetClaims(),
	}, nil
}
