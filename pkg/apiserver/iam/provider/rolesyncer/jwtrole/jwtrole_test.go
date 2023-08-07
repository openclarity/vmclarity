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

package jwtrole

import (
	"context"
	"testing"

	"github.com/openclarity/vmclarity/pkg/apiserver/iam"

	"github.com/spf13/viper"
	"gotest.tools/v3/assert"
)

func TestSync(t *testing.T) {
	t.Parallel()

	// Prepare
	ctx := context.Background()
	roleClaimKey := "role-claim-key"
	viper.Set(roleSyncerJwtRoleClaimEnvVar, roleClaimKey)

	// Create role syncer
	roleSyncer, err := New()
	assert.NilError(t, err)

	// Test
	for _, tt := range []struct {
		name   string
		user   *iam.User
		errMsg string
		roles  []string
	}{
		{
			name:   "nil user",
			user:   nil,
			errMsg: "no user provider",
		},
		{
			name:   "nil jwt claims",
			user:   &iam.User{},
			errMsg: "no user jwt claims found",
		},
		{
			name:   "empty jwt claims",
			user:   &iam.User{JwtClaims: map[string]interface{}{}},
			errMsg: "no user role data",
		},
		{
			name: "invalid claim type",
			user: &iam.User{
				JwtClaims: map[string]interface{}{
					roleClaimKey: "role",
				},
			},
			errMsg: "cannot extract roles",
		},
		{
			name: "role claim map",
			user: &iam.User{
				JwtClaims: map[string]interface{}{
					roleClaimKey: map[string]interface{}{"role": true},
				},
			},
			roles: []string{"role"},
		},
		{
			name: "role claim slice",
			user: &iam.User{
				JwtClaims: map[string]interface{}{
					roleClaimKey: []string{"role"},
				},
			},
			roles: []string{"role"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			syncErr := roleSyncer.Sync(ctx, tt.user)
			if tt.errMsg != "" {
				assert.ErrorContains(t, syncErr, tt.errMsg)
			} else {
				assert.NilError(t, syncErr)
				assert.DeepEqual(t, tt.user.GetRoles(), tt.roles)
			}
		})
	}
}
