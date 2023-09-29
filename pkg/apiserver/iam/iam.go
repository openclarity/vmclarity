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
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/authn/oidc"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"
)

const userCtxKey = "user"

// GetUserFromContext returns types.User from context.
func GetUserFromContext(ctx echo.Context) *types.User {
	ctxData := ctx.Get(userCtxKey)
	if ctxData == nil {
		return nil
	}
	user, _ := ctxData.(*types.User)
	return user
}

// NewMiddleware creates an OpenAPI middleware filter which handles request
// authentication and authorization using a specific Provider. Logic flow:
//
// request -> |Authenticator.Authenticate| -> |RoleSyncer.Sync| -> |Authorizer.CanPerform| -> success
//
// Provider will first authenticate the client from request data, synchronize
// user roles on success, and finally try to authorize the request. If
// successful, User will be available in context.
func NewMiddleware(service types.IAMService) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		if service == nil {
			return nil
		}

		// Remove user from request context
		eCtx := middleware.GetEchoContext(ctx)
		eCtx.Set(userCtxKey, nil)

		// Authenticate
		user, err := service.Authenticator().Authenticate(ctx, input.RequestValidationInput.Request)
		if err != nil {
			return fmt.Errorf("failed to authenticate request: %w", err)
		}

		// Authorize
		// Route permissions are defined as "asset:action", so we extract
		// the asset and action from the requested permission scope
		//for _, scope := range input.Scopes {
		//	// Fetch authorization request data
		//	reqScope := strings.Split(scope, ":")
		//	if len(reqScope) != 2 {
		//		return fmt.Errorf("unknown api asset:action found, got %s", scope)
		//	}
		//	asset, action := reqScope[0], reqScope[1]
		//
		//	// Authorize request
		//	authorized, err := service.Authorizer().CanPerform(ctx, user, asset, action)
		//	if err != nil {
		//		return fmt.Errorf("failed to check authorization: %w", err)
		//	}
		//	if !authorized {
		//		return fmt.Errorf("not allowed, missing permission to perform %s:%s", asset, action)
		//	}
		//}

		// Update request context with user data
		eCtx.Set(userCtxKey, user)

		return nil
	}
}

func NewService() (types.IAMService, error) {
	authenticator, err := oidc.New()
	if err != nil {
		return nil, err
	}
	return &iamService{
		authn: authenticator,
	}, nil
}

type iamService struct {
	authn types.Authenticator
	authz types.Authorizer
}

func (i iamService) Authenticator() types.Authenticator { return i.authn }

func (i iamService) Authorizer() types.Authorizer { return i.authz }
