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
	"net/http"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

const userCtxKey = "user"

// Logic flow:
// - Provider -> Interacts with IDP to authenticate, authorizer to authorize HTTP request, and syncer to sync user roles.
// - RoleSyncer -> Interacts with some kind of store (e.g. database or from JWT token claim) to fetch and sync user roles.
// - Authorizer -> Decides if a User can perform a given action on an asset based on provided rules.

// User defines an authenticated user. This object should never be created. It is
// only returned by IAM providers and GetUserFromContext.
type User struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`

	// JwtClaims helps holds JWT Claims
	JwtClaims map[string]interface{} `json:"jwt-claims"`
}

// GetUserFromContext returns User from context.
func GetUserFromContext(ctx echo.Context) *User {
	ctxData := ctx.Get(userCtxKey)
	if ctxData == nil {
		return nil
	}

	user, _ := ctxData.(*User)
	return user
}

// Injector implements client-side request authentication data injection.
type Injector interface {
	// Inject fetches the authentication data from a source and injects them into
	// request data (e.g. into Authorization header).
	Inject(ctx context.Context, request *http.Request) error
}

type RoleSyncerType string

// RoleSyncer implements server-side User role synchronization from a specific source.
type RoleSyncer interface {
	Type() RoleSyncerType
	// Sync synchronizes user roles either from a local source (e.g. JWT claims) or using an API request.
	Sync(ctx context.Context, user *User) error
}

// Authorizer implements authorization methods from a specific source.
type Authorizer interface {
	// CanPerform checks if User is allowed to perform an action on an asset.
	CanPerform(ctx context.Context, user *User, asset, action string) (bool, error)
}

// Provider implements server-side IAM synchronization policy.
type Provider interface {
	// RoleSyncer returns selected RoleSyncer to use for Provider.
	RoleSyncer() RoleSyncer

	// Authorizer returns selected Authorizer to use for Provider.
	Authorizer() Authorizer

	// Authenticate validates and verifies user auth details from request against an
	// auth provider. Only User.ID should be set and dependency fields required to
	// interact with RoleSyncer.
	Authenticate(ctx context.Context, request *http.Request) (*User, error)
}

// OapiFilterForProvider creates an OpenAPI filter function which handles request
// authentication and authorization using a specific Provider.
//
// Provider will first authenticate the client from request data, synchronize
// user roles on success, and finally try to authorize the request. This ensures
// that the User has all the required data to interact with IAM policies.
func OapiFilterForProvider(provider Provider) openapi3filter.AuthenticationFunc {
	if provider == nil {
		return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
			return nil
		}
	}

	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// TODO: Explore caching options to reduce checks against identity server

		// Authenticate
		user, err := provider.Authenticate(ctx, input.RequestValidationInput.Request)
		if err != nil {
			return fmt.Errorf("failed to authenticate user: %w", err)
		}

		// Sync user roles
		err = provider.RoleSyncer().Sync(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to sync user roles: %w", err)
		}

		// Authorize
		// Route permissions are defined as "asset:action", so we
		// extract the asset and action from the requested permission scope
		for _, scope := range input.Scopes {
			// Fetch authorization request data
			reqScope := strings.Split(scope, ":")
			if len(reqScope) != 2 {
				return fmt.Errorf("unknown asset:action found: %s", scope)
			}
			asset, action := reqScope[0], reqScope[1]

			// Authorize request
			authorized, err := provider.Authorizer().CanPerform(ctx, user, asset, action)
			if err != nil {
				return fmt.Errorf("failed to check authorization: %w", err)
			}
			if !authorized {
				return fmt.Errorf("not allowed, missing required permissions")
			}
		}

		// Success
		// Update request context with user data
		if eCtx := middleware.GetEchoContext(ctx); eCtx != nil {
			eCtx.Set(userCtxKey, user)
		}
		return nil
	}
}
