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

	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	user "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"
)

// New creates an authenticator which intercepts requests and checks for a
// correct Bearer token using OAuth2 introspection by sending the token to the
// introspection endpoint. On success, returns an iam.User with configured
// JwtClaims.
func New() (types.Authenticator, error) {
	// Load config
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config for Authenticator=OIDC: %w", err)
	}
	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config for Authenticator=OIDC: %w", err)
	}

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
	}, nil
}

type authService struct {
	resourceServer rs.ResourceServer
	mgmtClient     *management.Client
}

func (auth *authService) Authenticate(ctx context.Context, token string) (*types.User, error) {
	// Verify token against introspection endpoint
	jwtToken, err := rs.Introspect(ctx, auth.resourceServer, token)
	if err != nil || !jwtToken.IsActive() {
		return nil, fmt.Errorf("payload token is invalid")
	}

	// Load roles for given user
	roleData := jwtToken.GetClaim("roles")
	roles, _ := roleData.([]string)

	// Return user
	return &types.User{
		ID:    jwtToken.GetSubject(),
		Roles: roles,
	}, nil
}

func (auth *authService) CreateSA(ctx context.Context, username string) (string, error) {
	resp, err := auth.mgmtClient.AddMachineUser(ctx, &pb.AddMachineUserRequest{
		UserName:        username,
		Name:            username,
		Description:     "Service Account",
		AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	})
	return resp.UserId, err
}

func (auth *authService) DeleteSA(ctx context.Context, username string) error {
	resp, err := auth.mgmtClient.GetUserByLoginNameGlobal(ctx, &pb.GetUserByLoginNameGlobalRequest{
		LoginName: username,
	})
	if err != nil {
		return err
	}
	_, err = auth.mgmtClient.RemoveUser(ctx, &pb.RemoveUserRequest{Id: resp.User.Id})
	return err
}
