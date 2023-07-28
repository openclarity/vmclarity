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

package authorizer

import (
	"github.com/openclarity/vmclarity/backend/pkg/config"
	"github.com/openclarity/vmclarity/backend/pkg/iam"
)

// NewAuthorizer creates a new iam.Authorizer from config.
// TODO: Use Factory pattern when this supports multiple iam.Authorizer.
func NewAuthorizer(config config.Authorization) (iam.Authorizer, error) {
	return newLocalRBACAuthorizer(config.RBACLocal.RuleFilePath)
}
