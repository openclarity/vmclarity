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

package state

import (
	_ "embed"
	"github.com/openclarity/vmclarity/scanner/families"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"

	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
	"github.com/openclarity/vmclarity/scanner/config"
	"github.com/openclarity/vmclarity/scanner/families/sbom"
	"github.com/openclarity/vmclarity/scanner/families/types"
)

//go:embed testdata/effective-config.json
var effectiveScanConfigJSON string

func Test_appendEffectiveScanConfigAnnotation(t *testing.T) {
	type args struct {
		annotations *apitypes.Annotations
		config      *families.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *apitypes.Annotations
		wantErr bool
	}{
		{
			name: "annotations is nil",
			args: args{
				annotations: nil,
				config: &families.Config{
					SBOM: sbom.Config{
						Enabled:       true,
						AnalyzersList: []string{"syft"},
						Inputs: []types.Input{
							{
								Input:     "test",
								InputType: "dir",
							},
						},
						AnalyzersConfig: &config.Config{
							Analyzer: &config.Analyzer{
								AnalyzerList: []string{"syft"},
							},
						},
					},
				},
			},
			want: &apitypes.Annotations{
				{
					Key:   to.Ptr(effectiveScanConfigAnnotationKey),
					Value: to.Ptr(effectiveScanConfigJSON),
				},
			},
		},
		{
			name: "annotations is empty list",
			args: args{
				annotations: &apitypes.Annotations{},
				config: &families.Config{
					SBOM: sbom.Config{
						Enabled:       true,
						AnalyzersList: []string{"syft"},
						Inputs: []types.Input{
							{
								Input:     "test",
								InputType: "dir",
							},
						},
						AnalyzersConfig: &config.Config{
							Analyzer: &config.Analyzer{
								AnalyzerList: []string{"syft"},
							},
						},
					},
				},
			},
			want: &apitypes.Annotations{
				{
					Key:   to.Ptr(effectiveScanConfigAnnotationKey),
					Value: to.Ptr(effectiveScanConfigJSON),
				},
			},
		},
		{
			name: "annotations is not empty",
			args: args{
				annotations: &apitypes.Annotations{
					{
						Key:   to.Ptr("test"),
						Value: to.Ptr("test"),
					},
				},
				config: &families.Config{
					SBOM: sbom.Config{
						Enabled:       true,
						AnalyzersList: []string{"syft"},
						Inputs: []types.Input{
							{
								Input:     "test",
								InputType: "dir",
							},
						},
						AnalyzersConfig: &config.Config{
							Analyzer: &config.Analyzer{
								AnalyzerList: []string{"syft"},
							},
						},
					},
				},
			},
			want: &apitypes.Annotations{
				{
					Key:   to.Ptr("test"),
					Value: to.Ptr("test"),
				},
				{
					Key:   to.Ptr(effectiveScanConfigAnnotationKey),
					Value: to.Ptr(effectiveScanConfigJSON),
				},
			},
		},
		{
			name: "annotations is not empty and contains effective scan config, overwrite it",
			args: args{
				annotations: &apitypes.Annotations{
					{
						Key:   to.Ptr(effectiveScanConfigAnnotationKey),
						Value: to.Ptr("test"),
					},
				},
				config: &families.Config{
					SBOM: sbom.Config{
						Enabled:       true,
						AnalyzersList: []string{"syft"},
						Inputs: []types.Input{
							{
								Input:     "test",
								InputType: "dir",
							},
						},
						AnalyzersConfig: &config.Config{
							Analyzer: &config.Analyzer{
								AnalyzerList: []string{"syft"},
							},
						},
					},
				},
			},
			want: &apitypes.Annotations{
				{
					Key:   to.Ptr(effectiveScanConfigAnnotationKey),
					Value: to.Ptr(effectiveScanConfigJSON),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			got, err := appendEffectiveScanConfigAnnotation(tt.args.annotations, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("appendEffectiveScanConfigAnnotation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, w := range *tt.want {
				if *w.Key == effectiveScanConfigAnnotationKey {
					// In the case of effective scan config annotation the value won't
					// match because it is a formatted JSON.
					// We cannot use MatchJSON here because it doesn't work on pointers,
					// and the types.Annotations value is a pointer to string.
					g.Expect(*got).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Key": Equal(w.Key),
					})))
				} else {
					g.Expect(*got).To(ContainElement(gstruct.MatchAllFields(gstruct.Fields{
						"Key":   Equal(w.Key),
						"Value": Equal(w.Value),
					})))
				}
			}

			// Check the effective scan config JSON content
			for _, actual := range *got {
				g.Expect(actual).Should(HaveExistingField("Key"))
				g.Expect(actual).Should(HaveExistingField("Value"))
				// In the case of effective scan config annotation we check the JSON content of the value
				if *actual.Key == effectiveScanConfigAnnotationKey {
					g.Expect(*actual.Value).Should(MatchJSON(effectiveScanConfigJSON))
				}
			}
		})
	}
}
