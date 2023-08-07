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

package injector

import (
	"fmt"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/injector/envbearer"
)

// Options defines parameters for different iam.Injector.
//
// TODO: Extend to add support for other iam.Injector.
type Options struct {
	BearerTokenEnv string `json:"bearer-token-env"`
}

// New creates a new iam.Injector.
func New(kind models.IamInjector, options Options) (iam.Injector, error) {
	switch kind {
	case models.InjectorBearerToken:
		return envbearer.New(options.BearerTokenEnv)
	default:
		return nil, fmt.Errorf("factory not implemented for Injector=%s", kind)
	}
}
