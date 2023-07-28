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

package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zitadel/oidc/pkg/oidc"

	"github.com/zitadel/oidc/pkg/client/rs"

	"github.com/openclarity/vmclarity/pkg/apiserver/config"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/rolesyncer"
)

type oidcIDP struct {
	roleSyncer     iam.RoleSyncer
	authorizer     iam.Authorizer
	resourceServer rs.ResourceServer
}

// NewOIDCProvider creates a Provider which intercepts requests and checks for a
// correct Bearer token using OAuth2 introspection by sending the token to the
// introspection endpoint.
//
// TODO: Test against different OIDCs to check if this works. Tested against: Zitadel.
func newOIDCIdentityProvider(config config.AuthenticationOIDC, roleSyncer iam.RoleSyncer, authorizer iam.Authorizer) (iam.Provider, error) {
	// Check RoleSyncer support
	switch roleSyncerType := roleSyncer.Type(); roleSyncerType {
	case rolesyncer.RoleSyncerTypeJwt:
	// supported
	default:
		return nil, fmt.Errorf("unsupported role syncer type provided: %s", roleSyncerType)
	}

	// Add custom OIDC options
	var options []rs.Option
	if config.TokenURL != "" && config.IntrospectURL != "" {
		options = append(options, rs.WithStaticEndpoints("tokenUrl", "introspectUrl"))
	} else if config.TokenURL != "" || config.IntrospectURL != "" {
		return nil, fmt.Errorf("requires both OIDC token and introspect endpoints")
	}

	// Create resource server which provides introspection functionality
	resourceServer, err := rs.NewResourceServerClientCredentials(config.Issuer, config.ClientID, config.ClientSecret, options...)
	if err != nil {
		return nil, fmt.Errorf("could not create OIDC resource server: %w", err)
	}

	// Return OIDC Provider
	return &oidcIDP{
		roleSyncer:     roleSyncer,
		authorizer:     authorizer,
		resourceServer: resourceServer,
	}, nil
}

func (provider *oidcIDP) RoleSyncer() iam.RoleSyncer {
	return provider.roleSyncer
}

func (provider *oidcIDP) Authorizer() iam.Authorizer {
	return provider.authorizer
}

func (provider *oidcIDP) Authenticate(ctx context.Context, request *http.Request) (*iam.User, error) {
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

	// Return user
	return &iam.User{
		ID:        token.GetSubject(),
		JwtClaims: token.GetClaims(),
	}, nil
}
