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

package provider

import (
	"fmt"

	"github.com/openclarity/vmclarity/pkg/apiserver/config"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/authorizer"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/rolesyncer"
)

// NewProvider creates a new iam.Provider from config.
// TODO: Use Factory pattern when this supports multiple iam.Provider.
func NewProvider(config config.Config) (iam.Provider, error) {
	roleSyncer, err := rolesyncer.NewRoleSyncer(config.AuthRoleSynchronization)
	if err != nil {
		return nil, fmt.Errorf("failed creating role syncer: %w", err)
	}

	authzer, err := authorizer.NewAuthorizer(config.Authorization)
	if err != nil {
		return nil, fmt.Errorf("failed creating authorizer: %w", err)
	}

	return newOIDCIdentityProvider(config.Authentication.OIDC, roleSyncer, authzer)
}
