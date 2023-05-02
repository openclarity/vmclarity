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

package scanner

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func Test_getInitScanStatusVulnerabilitiesStateFromEnabled(t *testing.T) {
	type args struct {
		config *models.VulnerabilitiesConfig
	}
	tests := []struct {
		name string
		args args
		want *models.TargetScanStateState
	}{
		{
			name: "enabled",
			args: args{
				config: &models.VulnerabilitiesConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: stateToPointer(models.INIT),
		},
		{
			name: "disabled",
			args: args{
				config: &models.VulnerabilitiesConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil enabled",
			args: args{
				config: &models.VulnerabilitiesConfig{
					Enabled: nil,
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			want: stateToPointer(models.NOTSCANNED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInitScanStatusVulnerabilitiesStateFromEnabled(tt.args.config)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInitScanStatusVulnerabilitiesStateFromEnabled() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getInitScanStatusSecretsStateFromEnabled(t *testing.T) {
	type args struct {
		config *models.SecretsConfig
	}
	tests := []struct {
		name string
		args args
		want *models.TargetScanStateState
	}{
		{
			name: "enabled",
			args: args{
				config: &models.SecretsConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: stateToPointer(models.INIT),
		},
		{
			name: "disabled",
			args: args{
				config: &models.SecretsConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil enabled",
			args: args{
				config: &models.SecretsConfig{
					Enabled: nil,
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			want: stateToPointer(models.NOTSCANNED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInitScanStatusSecretsStateFromEnabled(tt.args.config)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInitScanStatusSecretsStateFromEnabled() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getInitScanStatusSbomStateFromEnabled(t *testing.T) {
	type args struct {
		config *models.SBOMConfig
	}
	tests := []struct {
		name string
		args args
		want *models.TargetScanStateState
	}{
		{
			name: "enabled",
			args: args{
				config: &models.SBOMConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: stateToPointer(models.INIT),
		},
		{
			name: "disabled",
			args: args{
				config: &models.SBOMConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil enabled",
			args: args{
				config: &models.SBOMConfig{
					Enabled: nil,
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			want: stateToPointer(models.NOTSCANNED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInitScanStatusSbomStateFromEnabled(tt.args.config)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInitScanStatusSbomStateFromEnabled() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getInitScanStatusRootkitsStateFromEnabled(t *testing.T) {
	type args struct {
		config *models.RootkitsConfig
	}
	tests := []struct {
		name string
		args args
		want *models.TargetScanStateState
	}{
		{
			name: "enabled",
			args: args{
				config: &models.RootkitsConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: stateToPointer(models.INIT),
		},
		{
			name: "disabled",
			args: args{
				config: &models.RootkitsConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil enabled",
			args: args{
				config: &models.RootkitsConfig{
					Enabled: nil,
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			want: stateToPointer(models.NOTSCANNED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInitScanStatusRootkitsStateFromEnabled(tt.args.config)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInitScanStatusRootkitsStateFromEnabled() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getInitScanStatusMisconfigurationsStateFromEnabled(t *testing.T) {
	type args struct {
		config *models.MisconfigurationsConfig
	}
	tests := []struct {
		name string
		args args
		want *models.TargetScanStateState
	}{
		{
			name: "enabled",
			args: args{
				config: &models.MisconfigurationsConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: stateToPointer(models.INIT),
		},
		{
			name: "disabled",
			args: args{
				config: &models.MisconfigurationsConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil enabled",
			args: args{
				config: &models.MisconfigurationsConfig{
					Enabled: nil,
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			want: stateToPointer(models.NOTSCANNED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInitScanStatusMisconfigurationsStateFromEnabled(tt.args.config)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInitScanStatusMisconfigurationsStateFromEnabled() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getInitScanStatusMalwareStateFromEnabled(t *testing.T) {
	type args struct {
		config *models.MalwareConfig
	}
	tests := []struct {
		name string
		args args
		want *models.TargetScanStateState
	}{
		{
			name: "enabled",
			args: args{
				config: &models.MalwareConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: stateToPointer(models.INIT),
		},
		{
			name: "disabled",
			args: args{
				config: &models.MalwareConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil enabled",
			args: args{
				config: &models.MalwareConfig{
					Enabled: nil,
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			want: stateToPointer(models.NOTSCANNED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInitScanStatusMalwareStateFromEnabled(tt.args.config)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInitScanStatusMalwareStateFromEnabled() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_getInitScanStatusExploitsStateFromEnabled(t *testing.T) {
	type args struct {
		config *models.ExploitsConfig
	}
	tests := []struct {
		name string
		args args
		want *models.TargetScanStateState
	}{
		{
			name: "enabled",
			args: args{
				config: &models.ExploitsConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: stateToPointer(models.INIT),
		},
		{
			name: "disabled",
			args: args{
				config: &models.ExploitsConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil enabled",
			args: args{
				config: &models.ExploitsConfig{
					Enabled: nil,
				},
			},
			want: stateToPointer(models.NOTSCANNED),
		},
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			want: stateToPointer(models.NOTSCANNED),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInitScanStatusExploitsStateFromEnabled(tt.args.config)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInitScanStatusExploitsStateFromEnabled() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
