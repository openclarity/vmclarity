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

package authn

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"strings"

	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"

	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/zitadel/oidc/pkg/client/rs"
)

type oidcAuth struct {
	config         *oauth2.Config
	provider       *oidc.Provider
	resourceServer rs.ResourceServer
	isZitadel      bool
}

// New creates an authenticator which intercepts requests and validates the
// Bearer token via OIDC introspection.
func New() (types.Authenticator, error) {
	config := LoadConfig()

	// Create resource server which provides introspection functionality
	resourceServer, err := rs.NewResourceServerClientCredentials(config.Issuer, config.ClientID, config.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("could not create resource server for Authenticator=OIDC: %w", err)
	}

	// Create provider
	provider, err := oidc.NewProvider(context.Background(), config.Issuer)
	if err != nil {
		return nil, err
	}

	// Return OIDC Authenticator
	return &oidcAuth{
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  config.RedirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "profile"},
		},
		provider:       provider,
		resourceServer: resourceServer,
		isZitadel:      config.UseZitadel,
	}, nil
}

func (auth *oidcAuth) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return auth.config.AuthCodeURL(state, opts...)
}

func (auth *oidcAuth) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return auth.config.Exchange(ctx, code, opts...)
}

func (auth *oidcAuth) Verify(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	return auth.provider.Verifier(&oidc.Config{
		ClientID: auth.config.ClientID,
	}).Verify(ctx, rawIDToken)
}

func (auth *oidcAuth) Introspect(ctx context.Context, req *http.Request) (*types.UserInfo, error) {
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
	// Return authenticated user info
	return &types.UserInfo{
		UserInfo: oidc.UserInfo{
			Subject:       jwtToken.GetSubject(),
			Profile:       jwtToken.GetProfile(),
			Email:         jwtToken.GetEmail(),
			EmailVerified: jwtToken.IsEmailVerified(),
		},
		FromGenericOIDC: !auth.isZitadel,
		FromZitadelOIDC: auth.isZitadel,
	}, nil
}

func extractToken(r *http.Request) (string, error) {
	auth := r.Header.Get("authorization")
	if auth == "" {
		// http.StatusUnauthorized
		return "", fmt.Errorf("auth header missing")
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		// http.StatusUnauthorized
		return "", fmt.Errorf("invalid auth header")
	}
	return strings.TrimPrefix(auth, "Bearer "), nil
}
