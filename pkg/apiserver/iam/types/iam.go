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
	"net/http"
)

type IAMService interface {
	Authenticator() Authenticator
	Authorizer() Authorizer
}

//	type ScansTable interface {
//		GetScans(params models.GetScansParams) (models.Scans, error)
//		GetScan(scanID models.ScanID, params models.GetScansScanIDParams) (models.Scan, error)
//
//		CreateScan(scan models.Scan) (models.Scan, error)
//		UpdateScan(scan models.Scan, params models.PatchScansScanIDParams) (models.Scan, error)
//		SaveScan(scan models.Scan, params models.PutScansScanIDParams) (models.Scan, error)
//
//		DeleteScan(scanID models.ScanID) error
//	}

// AuthUser defines an authenticated user.
type AuthUser struct {
	ID string

	// Data from auth providers
	FromOIDC *AuthFromOIDC
}

// AuthFromOIDC defines data from OIDC Authenticator
type AuthFromOIDC struct {
	Claims map[string]interface{}
}

// Authenticator defines authentication service.
type Authenticator interface {
	// Authenticate authenticates a request payload against authentication service.
	Authenticate(ctx context.Context, req *http.Request) (*AuthUser, error)
}

type User struct {
	ID      string
	Name    string
	Email   string
	Roles   []string
	Banned  bool
	Machine bool
}

type UserStore interface {
	CurrentUser(ctx context.Context) (*User, error)

	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, user *User) error

	CreateCredentials(ctx context.Context, user *User) error
	RevokeCredentials(ctx context.Context, user *User, token string) error
}

// Authorizer defines authorization service.
type Authorizer interface {
	// CanPerform checks if User is allowed to perform an action on an asset based on
	// some predefined rules.
	CanPerform(ctx context.Context, user interface{}, asset, action string) (bool, error)
}
