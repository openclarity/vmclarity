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
	"github.com/openclarity/vmclarity/scanner/common"
	sbomtypes "github.com/openclarity/vmclarity/scanner/families/sbom/types"
	secrettypes "github.com/openclarity/vmclarity/scanner/families/secrets/types"
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/core/to"
	malware "github.com/openclarity/vmclarity/scanner/families/malware/types"
	vulnerabilities "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/types"
)

func Test_SetMountPointsForFamiliesInput(t *testing.T) {
	type args struct {
		mountPoints    []string
		familiesConfig *Config
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "sbom, vuls, secrets and malware are enabled",
			args: args{
				mountPoints: []string{"/mnt/snapshot1"},
				familiesConfig: &Config{
					SBOM: sbomtypes.Config{
						Enabled: true,
						Inputs:  nil,
					},
					Vulnerabilities: vulnerabilities.Config{
						Enabled:       true,
						Inputs:        nil,
						InputFromSbom: false,
					},
					Secrets: secrettypes.Config{
						Enabled: true,
						Inputs:  nil,
					},
					Malware: malware.Config{
						Enabled: true,
						Inputs:  nil,
					},
				},
			},
			want: &Config{
				SBOM: sbomtypes.Config{
					Enabled: true,
					Inputs: []common.ScanInput{
						{
							Input:     "/mnt/snapshot1",
							InputType: common.ROOTFS,
						},
					},
				},
				Vulnerabilities: vulnerabilities.Config{
					Enabled:       true,
					InputFromSbom: true,
				},
				Secrets: secrettypes.Config{
					Enabled: true,
					Inputs: []common.ScanInput{
						{
							StripPathFromResult: to.Ptr(true),
							Input:               "/mnt/snapshot1",
							InputType:           common.ROOTFS,
						},
					},
				},
				Malware: malware.Config{
					Enabled: true,
					Inputs: []common.ScanInput{
						{
							StripPathFromResult: to.Ptr(true),
							Input:               "/mnt/snapshot1",
							InputType:           common.ROOTFS,
						},
					},
				},
			},
		},
		{
			name: "only vuls enabled",
			args: args{
				mountPoints: []string{"/mnt/snapshot1"},
				familiesConfig: &Config{
					Vulnerabilities: vulnerabilities.Config{
						Enabled:       true,
						Inputs:        nil,
						InputFromSbom: false,
					},
				},
			},
			want: &Config{
				Vulnerabilities: vulnerabilities.Config{
					Enabled: true,
					Inputs: []common.ScanInput{
						{
							Input:     "/mnt/snapshot1",
							InputType: common.ROOTFS,
						},
					},
					InputFromSbom: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.familiesConfig.AddInputs(common.ROOTFS, tt.args.mountPoints)
			got := tt.args.familiesConfig
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddInputs() = %v, want %v", got, tt.want)
			}
		})
	}
}
