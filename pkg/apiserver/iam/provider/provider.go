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
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/provider/authenticator"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/provider/authorizer"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/provider/rolesyncer"
)

// New creates a new iam.Provider.
//
// TODO: Add support for creating dynamic types.
func New() (iam.Provider, error) {
	auth, err := authenticator.New(models.AuthenticatorOIDC)
	if err != nil {
		return nil, fmt.Errorf("failed creating authenticator: %w", err)
	}

	roleSyncer, err := rolesyncer.New(models.RoleSyncerJWT)
	if err != nil {
		return nil, fmt.Errorf("failed creating rolesyncer: %w", err)
	}

	authz, err := authorizer.New(models.AuthorizerLocalRBAC)
	if err != nil {
		return nil, fmt.Errorf("failed creating authorizer: %w", err)
	}

	return &provider{
		authenticator: auth,
		roleSyncer:    roleSyncer,
		authorizer:    authz,
	}, nil
}

// NewFromParams creates a new iam.Provider from params.
func NewFromParams(authenticator iam.Authenticator, roleSyncer iam.RoleSyncer, authorizer iam.Authorizer) (iam.Provider, error) {
	if authenticator == nil {
		return nil, fmt.Errorf("cannot create Provider with nil Authenticator")
	}
	if roleSyncer == nil {
		return nil, fmt.Errorf("cannot create Provider with nil RoleSyncer")
	}
	if authorizer == nil {
		return nil, fmt.Errorf("cannot create Provider with nil Authorizer")
	}

	return &provider{
		authenticator: authenticator,
		roleSyncer:    roleSyncer,
		authorizer:    authorizer,
	}, nil
}

type provider struct {
	authenticator iam.Authenticator
	roleSyncer    iam.RoleSyncer
	authorizer    iam.Authorizer
}

func (p *provider) Authenticator() iam.Authenticator { return p.authenticator }

func (p *provider) RoleSyncer() iam.RoleSyncer { return p.roleSyncer }

func (p *provider) Authorizer() iam.Authorizer { return p.authorizer }
