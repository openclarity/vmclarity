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
	"net/http"
	"sync"
)

// User defines an authenticated user. This object should never be manually created.
type User struct {
	ID        string                 `json:"id"`         // ID holds unique UserID
	JwtClaims map[string]interface{} `json:"jwt-claims"` // JwtClaims holds JWT Claims data

	mu    sync.Mutex
	roles []string
}

func (user *User) GetRoles() []string {
	return user.roles
}

func (user *User) SetRoles(roles []string) {
	user.mu.Lock()
	defer user.mu.Unlock()
	user.roles = roles
}

// Injector implements client-side request authentication data injection. It is
// intended to be used by clients to allow adding auth data from various sources
// (e.g. from file, env variable, etc)
type Injector interface {
	// Inject fetches the authentication data from a source and injects them into
	// request data (e.g. token from env var into Authorization header).
	Inject(ctx context.Context, request *http.Request) error
}

// Authenticator implements server-side authentication methods. It is only
// intended to be used by server middleware to enable authentication.
type Authenticator interface {
	// Authenticate validates and verifies user auth details from request against an
	// auth provider. Publicly exported User fields should be set, or at least User.ID.
	Authenticate(ctx context.Context, request *http.Request) (*User, error)
}

// RoleSyncer implements server-side User role synchronization. It is intended to
// be used by server middleware to synchronize User roles. However, it can be
// used anywhere in the server stack for User role synchronization.
type RoleSyncer interface {
	// Sync synchronizes user roles either from a local source (e.g. JWT claims) or using an API request.
	Sync(ctx context.Context, user *User) error
}

// Authorizer implements server-side authorization methods. It is intended to be
// used by server middleware to check if User can perform an action on an asset.
// However, it can be used anywhere in the server stack for authorization purposes.
type Authorizer interface {
	// CanPerform checks if User is allowed to perform an action on an asset.
	CanPerform(ctx context.Context, user *User, asset, action string) (bool, error)
}

// Provider unifies server-side components to enable IAM operations.
type Provider interface {
	Authenticator() Authenticator
	RoleSyncer() RoleSyncer
	Authorizer() Authorizer
}
