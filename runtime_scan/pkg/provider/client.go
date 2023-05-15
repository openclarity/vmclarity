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
	"context"

	"github.com/openclarity/vmclarity/api/models"
)

type Client interface {
	// Kind returns provider type
	Kind() models.CloudProvider
	// DiscoverScopes returns a list of discovered scopes
	DiscoverScopes(ctx context.Context) (*models.Scopes, error)
	// DiscoverTargets returns list of Targets in scanScope
	DiscoverTargets(ctx context.Context, scanScope *models.ScanScopeType) ([]models.TargetType, error)
	RunTargetScan(context.Context, *ScanJobConfig) (models.TargetScanStateState, error)
	RemoveTargetScan(context.Context, *ScanJobConfig) (models.ResourceCleanupState, error)
}
