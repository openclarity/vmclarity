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

package findingkey

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"

	"github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/cli/pkg/utils"
)

func TestGenerateFindingKey(t *testing.T) {
	rootkitFindingInfo := types.RootkitFindingInfo{
		Message:     utils.PointerTo("Message"),
		RootkitName: utils.PointerTo("RootkitName"),
		RootkitType: utils.PointerTo(types.RootkitType("RootkitType")),
	}
	exploitFindingInfo := types.ExploitFindingInfo{
		CveID:       utils.PointerTo("CveID"),
		Description: utils.PointerTo("Description"),
		Name:        utils.PointerTo("Name"),
		SourceDB:    utils.PointerTo("SourceDB"),
		Title:       utils.PointerTo("Title"),
		Urls:        utils.PointerTo([]string{"url1", "url2"}),
	}
	vulFindingInfo := types.VulnerabilityFindingInfo{
		Package: &types.Package{
			Name:    utils.PointerTo("Package.Name"),
			Version: utils.PointerTo("Package.Version"),
		},
		VulnerabilityName: utils.PointerTo("VulnerabilityName"),
	}
	malwareFindingInfo := types.MalwareFindingInfo{
		MalwareName: utils.PointerTo("MalwareName"),
		MalwareType: utils.PointerTo("MalwareType"),
		Path:        utils.PointerTo("Path"),
		RuleName:    utils.PointerTo("RuleName"),
	}
	miscFindingInfo := types.MisconfigurationFindingInfo{
		Message:     utils.PointerTo("Message"),
		ScannerName: utils.PointerTo("ScannerName"),
		TestID:      utils.PointerTo("TestID"),
	}
	secretFindingInfo := types.SecretFindingInfo{
		EndColumn:   utils.PointerTo(1),
		Fingerprint: utils.PointerTo("Fingerprint"),
		StartColumn: utils.PointerTo(2),
	}
	pkgFindingInfo := types.PackageFindingInfo{
		Name:    utils.PointerTo("Name"),
		Version: utils.PointerTo("Version"),
	}

	type args struct {
		findingInfo *types.Finding_FindingInfo
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "exploit",
			args: args{
				findingInfo: createFindingInfo(t, exploitFindingInfo),
			},
			want:    GenerateExploitKey(exploitFindingInfo).ExploitString(),
			wantErr: false,
		},
		{
			name: "vul",
			args: args{
				findingInfo: createFindingInfo(t, vulFindingInfo),
			},
			want:    GenerateVulnerabilityKey(vulFindingInfo).VulnerabilityString(),
			wantErr: false,
		},
		{
			name: "malware",
			args: args{
				findingInfo: createFindingInfo(t, malwareFindingInfo),
			},
			want:    GenerateMalwareKey(malwareFindingInfo).MalwareString(),
			wantErr: false,
		},
		{
			name: "misc",
			args: args{
				findingInfo: createFindingInfo(t, miscFindingInfo),
			},
			want:    GenerateMisconfigurationKey(miscFindingInfo).MisconfigurationString(),
			wantErr: false,
		},
		{
			name: "rootkit",
			args: args{
				findingInfo: createFindingInfo(t, rootkitFindingInfo),
			},
			want:    GenerateRootkitKey(rootkitFindingInfo).RootkitString(),
			wantErr: false,
		},
		{
			name: "secret",
			args: args{
				findingInfo: createFindingInfo(t, secretFindingInfo),
			},
			want:    GenerateSecretKey(secretFindingInfo).SecretString(),
			wantErr: false,
		},
		{
			name: "pkg",
			args: args{
				findingInfo: createFindingInfo(t, pkgFindingInfo),
			},
			want:    GeneratePackageKey(pkgFindingInfo).PackageString(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateFindingKey(tt.args.findingInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFindingKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GenerateFindingKey() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func createFindingInfo(t *testing.T, info interface{}) *types.Finding_FindingInfo {
	t.Helper()
	var err error
	findingInfoB := types.Finding_FindingInfo{}
	switch fInfo := info.(type) {
	case types.RootkitFindingInfo:
		err = findingInfoB.FromRootkitFindingInfo(fInfo)
	case types.ExploitFindingInfo:
		err = findingInfoB.FromExploitFindingInfo(fInfo)
	case types.SecretFindingInfo:
		err = findingInfoB.FromSecretFindingInfo(fInfo)
	case types.MisconfigurationFindingInfo:
		err = findingInfoB.FromMisconfigurationFindingInfo(fInfo)
	case types.MalwareFindingInfo:
		err = findingInfoB.FromMalwareFindingInfo(fInfo)
	case types.VulnerabilityFindingInfo:
		err = findingInfoB.FromVulnerabilityFindingInfo(fInfo)
	case types.PackageFindingInfo:
		err = findingInfoB.FromPackageFindingInfo(fInfo)
	}
	assert.NilError(t, err)
	return &findingInfoB
}
