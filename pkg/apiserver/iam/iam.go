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
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

const userCtxKey = "user"

// GetUserFromContext returns User from context.
func GetUserFromContext(ctx echo.Context) *User {
	ctxData := ctx.Get(userCtxKey)
	if ctxData == nil {
		return nil
	}

	user, _ := ctxData.(*User)
	return user
}

// OapiFilterForProvider creates an OpenAPI middleware filter function which
// handles request authentication and authorization using a specific Provider.
// Logic flow:
//
// request -> |Authenticator.Authenticate| -> |RoleSyncer.Sync| -> |Authorizer.CanPerform| -> success
//
// Provider will first authenticate the client from request data, synchronize
// user roles on success, and finally try to authorize the request. If
// successful, User will be available for fetching in context.
func OapiFilterForProvider(provider Provider) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// TODO: Explore caching options to reduce checks against identity server

		// Authenticate
		user, err := provider.Authenticator().Authenticate(ctx, input.RequestValidationInput.Request)
		if err != nil {
			return fmt.Errorf("failed to authenticate user: %w", err)
		}

		// Sync user roles
		err = provider.RoleSyncer().Sync(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to sync user roles: %w", err)
		}

		// Authorize
		// Route permissions are defined as "asset:action", so we extract
		// the asset and action from the requested permission scope
		for _, scope := range input.Scopes {
			// Fetch authorization request data
			reqScope := strings.Split(scope, ":")
			if len(reqScope) != 2 {
				return fmt.Errorf("unknown api asset:action found, got %s", scope)
			}
			asset, action := reqScope[0], reqScope[1]

			// Authorize request
			authorized, err := provider.Authorizer().CanPerform(ctx, user, asset, action)
			if err != nil {
				return fmt.Errorf("failed to check authorization: %w", err)
			}
			if !authorized {
				return fmt.Errorf("not allowed, missing permission to perform %s:%s", asset, action)
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
