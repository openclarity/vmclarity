// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/yudai/gojsondiff/formatter"
	"gotest.tools/assert"
)

func Test_handleVulnerabilityWithExistingKey(t *testing.T) {
	highVul := Vulnerability{
		ID:       "highVul",
		Severity: "HIGH",
		Package: Package{
			Name:    "pkg-name",
			Version: "pkg-version",
		},
	}
	type args struct {
		vulnerability      Vulnerability
		otherVulnerability Vulnerability
	}
	tests := []struct {
		name string
		args args
		want Vulnerability
	}{
		{
			name: "identical vulnerability",
			args: args{
				vulnerability:      highVul,
				otherVulnerability: highVul,
			},
			want: highVul,
		},
		{
			name: "different fix versions",
			args: args{
				vulnerability: Vulnerability{
					ID: "highVul",
					Fix: Fix{
						Versions: []string{"1", "3"},
					},
				},
				otherVulnerability: Vulnerability{
					ID: "highVul",
					Fix: Fix{
						Versions: []string{"1", "2"},
					},
				},
			},
			want: Vulnerability{
				ID: "highVul",
				Fix: Fix{
					Versions: []string{"1", "2", "3"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleVulnerabilityWithExistingKey(tt.args.vulnerability, tt.args.otherVulnerability)
			assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreTypes(VulnerabilityDiff{}.ASCIIDiff))
		})
	}
}

func Test_getDiff(t *testing.T) {
	type args struct {
		vulnerability          Vulnerability
		compareToVulnerability Vulnerability
		compareToID            string
	}
	tests := []struct {
		name    string
		args    args
		want    *VulnerabilityDiff
		wantErr bool
	}{
		{
			name: "diff in fix",
			args: args{
				vulnerability: Vulnerability{
					ID: "id",
					Fix: Fix{
						Versions: []string{"1", "3"},
						State:    "not fixed",
					},
				},
				compareToVulnerability: Vulnerability{
					ID: "id",
					Fix: Fix{
						Versions: []string{"1", "2"},
						State:    "fixed",
					},
				},
				compareToID: "compareToID",
			},
			want: &VulnerabilityDiff{
				CompareToID: "compareToID",
				JSONDiff: map[string]interface{}{
					"fix": map[string]interface{}{
						"state": []interface{}{"fixed", "not fixed"},
						"versions": map[string]interface{}{
							"1":  []interface{}{"2", "3"}, // diff in array index 1
							"_t": "a",                     // sign for delta json
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "diff in links",
			args: args{
				vulnerability: Vulnerability{
					ID:    "id",
					Links: []string{"link1", "link2"},
				},
				compareToVulnerability: Vulnerability{
					ID:    "id",
					Links: []string{"link1", "link3", "link4"},
				},
				compareToID: "compareToID",
			},
			want: &VulnerabilityDiff{
				CompareToID: "compareToID",
				JSONDiff: map[string]interface{}{
					"links": map[string]interface{}{
						"1": []interface{}{"link4", "link2"},
						// "_1" means object was deleted from index 1
						// more info in github.com/yudai/gojsondiff@v1.0.0/formatter/delta.go
						"_1": []interface{}{"link3", 0, formatter.DeltaDelete},
						"_t": "a", // sign for delta json
					},
				},
				// ASCIIDiff: "{\n   \"cvss\": null,\n   \"distro\": {\n     \"idLike\": null,\n     \"name\": \"\",\n     \"version\": \"\"\n   },\n   \"fix\": {\n     \"state\": \"\",\n     \"versions\": null\n   },\n   \"id\": \"id\",\n   \"layerID\": \"\",\n   \"links\": [\n     0: \"link1\",\n-    1: \"link4\",\n+    1: \"link2\",\n-    1: \"link3\"\n     2: \"link4\"\n   ],\n   \"package\": {\n     \"cpes\": null,\n     \"language\": \"\",\n     \"licenses\": null,\n     \"name\": \"\",\n     \"purl\": \"\",\n     \"type\": \"\",\n     \"version\": \"\"\n   },\n   \"path\": \"\"\n }\n        ",
			},
			wantErr: false,
		},
		{
			name: "no diff - CVSS sort is needed",
			args: args{
				vulnerability: Vulnerability{
					CVSS: []CVSS{
						{
							Version: "3",
							Vector:  "456",
						},
						{
							Version: "2",
							Vector:  "123",
						},
					},
				},
				compareToVulnerability: Vulnerability{
					CVSS: []CVSS{
						{
							Version:        "2",
							Vector:         "123",
							Metrics:        CvssMetrics{},
							VendorMetadata: nil,
						},
						{
							Version:        "3",
							Vector:         "456",
							Metrics:        CvssMetrics{},
							VendorMetadata: nil,
						},
					},
				},
				compareToID: "compareToID",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDiff(tt.args.vulnerability, tt.args.compareToVulnerability, tt.args.compareToID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreTypes(VulnerabilityDiff{}.ASCIIDiff))
		})
	}
}

func Test_sortArrays(t *testing.T) {
	type args struct {
		vulnerability Vulnerability
	}
	tests := []struct {
		name string
		args args
		want Vulnerability
	}{
		{
			name: "sort",
			args: args{
				vulnerability: Vulnerability{
					Links: []string{"link2", "link1"},
					CVSS: []CVSS{
						{
							Version: "3",
							Vector:  "456",
						},
						{
							Version: "2",
							Vector:  "123",
						},
					},
					Fix: Fix{
						Versions: []string{"ver2", "ver1"},
					},
					Package: Package{
						Licenses: []string{"lic2", "lic1"},
						CPEs:     []string{"cpes2", "cpes1"},
					},
				},
			},
			want: Vulnerability{
				Links: []string{"link1", "link2"},
				CVSS: []CVSS{
					{
						Version: "2",
						Vector:  "123",
					},
					{
						Version: "3",
						Vector:  "456",
					},
				},
				Fix: Fix{
					Versions: []string{"ver1", "ver2"},
				},
				Package: Package{
					Licenses: []string{"lic1", "lic2"},
					CPEs:     []string{"cpes1", "cpes2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.vulnerability.sorted(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortArrays() = %v, want %v", got, tt.want)
			}
		})
	}
}
