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

// TODO: Implement via https://auth0.com/docs/quickstart/webapp/golang/interactive

package types

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/api/models"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrAlreadyExists = fmt.Errorf("already exists")
)

// UserInfo defines an authenticated (OIDC) user.
type UserInfo struct {
	oidc.UserInfo
	oidc.Provider
	// Data to indicate auth source
	FromGenericOIDC bool
	FromZitadelOIDC bool
}

// Authenticator defines (OIDC) authentication service.
type Authenticator interface {
	// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page that asks for
	// permissions for the required scopes explicitly.
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	// Exchange converts an authorization code into a token.
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	// Verify verifies that an *oauth2.Token is a valid *oidc.IDToken.
	Verify(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error)
	// Introspect fetches UserInfo data from OIDC IDP using introspect API. Consider caching.
	Introspect(ctx context.Context, req *http.Request) (*UserInfo, error)
}

// Authorizer defines authorization service.
type Authorizer interface {
	// CanPerform checks if user is allowed to perform an action on an asset based on specs.
	CanPerform(user models.User, category, action, asset string) (bool, error)
}

// AuthStore defines a (db) store to interact with user and auth data.
type AuthStore interface {
	GetUserFromInfo(info *UserInfo) (models.User, error)

	GetUsers(params models.GetUsersParams) (models.Users, error)
	GetUser(userID models.UserID) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(userID models.UserID) error

	GetUserAuth(userID models.UserID) (models.UserAuths, error)
	CreateUserAuth(userID models.UserID, credType models.CredentialType, credExpiry *models.CredentialExpiry) (models.UserCred, error)
	RevokeUserAuth(userID models.UserID, userAuth models.UserAuth) error
}
