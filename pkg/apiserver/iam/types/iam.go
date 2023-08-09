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

package types

import (
	"context"
)

// User defines an authenticated user.
type User struct {
	ID    string   `json:"ID"`
	Roles []string `json:"roles"`
}

// Authenticator defines authentication service.
type Authenticator interface {
	// Authenticate authenticates a request payload against authentication service.
	Authenticate(ctx context.Context, token string) (*User, error)

	CreateSA(ctx context.Context, username string) (string, error)
	DeleteSA(ctx context.Context, username string) error
}

// Authorizer defines authorization service.
type Authorizer interface {
	// CanPerform checks if User is allowed to perform an action on an asset based on
	// some predefined rules.
	CanPerform(ctx context.Context, user *User, asset, action string) (bool, error)
}
