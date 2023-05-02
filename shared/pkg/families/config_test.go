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

package families

import (
	"testing"
	"time"

	"github.com/anchore/syft/syft/source"
	"github.com/google/go-cmp/cmp"
	kubeclarityConfig "github.com/openclarity/kubeclarity/shared/pkg/config"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/families/malware"
	malwareconfig "github.com/openclarity/vmclarity/shared/pkg/families/malware/clam/config"
	malwarecommon "github.com/openclarity/vmclarity/shared/pkg/families/malware/common"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	secretscommon "github.com/openclarity/vmclarity/shared/pkg/families/secrets/common"
	gitleaksconfig "github.com/openclarity/vmclarity/shared/pkg/families/secrets/gitleaks/config"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func Test_userSBOMConfigToFamiliesSbomConfig(t *testing.T) {
	type args struct {
		sbomConfig *models.SBOMConfig
	}
	type returns struct {
		config sbom.Config
	}
	tests := []struct {
		name string
		args args
		want returns
	}{
		{
			name: "No SBOM Config",
			args: args{
				sbomConfig: nil,
			},
			want: returns{
				config: sbom.Config{},
			},
		},
		{
			name: "Missing Enabled",
			args: args{
				sbomConfig: &models.SBOMConfig{},
			},
			want: returns{
				config: sbom.Config{},
			},
		},
		{
			name: "Disabled",
			args: args{
				sbomConfig: &models.SBOMConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: returns{
				config: sbom.Config{},
			},
		},
		{
			name: "Enabled",
			args: args{
				sbomConfig: &models.SBOMConfig{
					Enabled: utils.BoolPtr(true),
				},
			},
			want: returns{
				config: sbom.Config{
					Enabled:       true,
					AnalyzersList: []string{"syft", "trivy"},
					AnalyzersConfig: &kubeclarityConfig.Config{
						Registry: &kubeclarityConfig.Registry{},
						Analyzer: &kubeclarityConfig.Analyzer{
							OutputFormat: "cyclonedx",
							TrivyConfig: kubeclarityConfig.AnalyzerTrivyConfig{
								Timeout: TrivyTimeout,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userSBOMConfigToFamiliesSbomConfig(tt.args.sbomConfig)
			if diff := cmp.Diff(tt.want.config, got); diff != "" {
				t.Errorf("userSBOMConfigToFamiliesSbomConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_userVulnConfigToFamiliesVulnConfig(t *testing.T) {
	type args struct {
		vulnerabilitiesConfig *models.VulnerabilitiesConfig
		trivyServerAddress    string
		grypeServerAddress    string
	}
	type returns struct {
		config vulnerabilities.Config
	}
	tests := []struct {
		name string
		args args
		want returns
	}{
		{
			name: "No Vulnerability Config",
			args: args{
				vulnerabilitiesConfig: nil,
				trivyServerAddress:    "http://10.0.0.1:9992",
				grypeServerAddress:    "",
			},
			want: returns{
				config: vulnerabilities.Config{},
			},
		},
		{
			name: "Missing Enabled",
			args: args{
				vulnerabilitiesConfig: &models.VulnerabilitiesConfig{},
				trivyServerAddress:    "http://10.0.0.1:9992",
				grypeServerAddress:    "",
			},
			want: returns{
				config: vulnerabilities.Config{},
			},
		},
		{
			name: "Disabled",
			args: args{
				vulnerabilitiesConfig: &models.VulnerabilitiesConfig{
					Enabled: utils.BoolPtr(false),
				},
				trivyServerAddress: "http://10.0.0.1:9992",
				grypeServerAddress: "",
			},
			want: returns{
				config: vulnerabilities.Config{},
			},
		},
		{
			name: "Enabled",
			args: args{
				vulnerabilitiesConfig: &models.VulnerabilitiesConfig{
					Enabled: utils.BoolPtr(true),
				},
				trivyServerAddress: "http://10.0.0.1:9992",
				grypeServerAddress: "",
			},
			want: returns{
				config: vulnerabilities.Config{
					Enabled: true,
					// TODO(sambetts) This choice should come from the user's configuration
					ScannersList: []string{"grype", "trivy"},
					ScannersConfig: &kubeclarityConfig.Config{
						// TODO(sambetts) The user needs to be able to provide this configuration
						Registry: &kubeclarityConfig.Registry{},
						Scanner: &kubeclarityConfig.Scanner{
							GrypeConfig: kubeclarityConfig.GrypeConfig{
								Mode: kubeclarityConfig.ModeLocal,
								LocalGrypeConfig: kubeclarityConfig.LocalGrypeConfig{
									UpdateDB:   true,
									DBRootDir:  "/tmp/",
									ListingURL: "https://toolbox-data.anchore.io/grype/databases/listing.json",
									Scope:      source.SquashedScope,
								},
							},
							TrivyConfig: kubeclarityConfig.ScannerTrivyConfig{
								Timeout:    TrivyTimeout,
								ServerAddr: "http://10.0.0.1:9992",
							},
						},
					},
				},
			},
		},
		{
			name: "Enabled with grype server",
			args: args{
				vulnerabilitiesConfig: &models.VulnerabilitiesConfig{
					Enabled: utils.BoolPtr(true),
				},
				trivyServerAddress: "http://10.0.0.1:9992",
				grypeServerAddress: "10.0.0.1:9991",
			},
			want: returns{
				config: vulnerabilities.Config{
					Enabled: true,
					// TODO(sambetts) This choice should come from the user's configuration
					ScannersList: []string{"grype", "trivy"},
					ScannersConfig: &kubeclarityConfig.Config{
						// TODO(sambetts) The user needs to be able to provide this configuration
						Registry: &kubeclarityConfig.Registry{},
						Scanner: &kubeclarityConfig.Scanner{
							GrypeConfig: kubeclarityConfig.GrypeConfig{
								Mode: kubeclarityConfig.ModeRemote,
								RemoteGrypeConfig: kubeclarityConfig.RemoteGrypeConfig{
									GrypeServerAddress: "10.0.0.1:9991",
									GrypeServerTimeout: 2 * time.Minute,
								},
							},
							TrivyConfig: kubeclarityConfig.ScannerTrivyConfig{
								Timeout:    TrivyTimeout,
								ServerAddr: "http://10.0.0.1:9992",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userVulnConfigToFamiliesVulnConfig(tt.args.vulnerabilitiesConfig, tt.args.trivyServerAddress, tt.args.grypeServerAddress)
			if diff := cmp.Diff(tt.want.config, got); diff != "" {
				t.Errorf("userVulnConfigToFamiliesVulnConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_userSecretsConfigToFamiliesSecretsConfig(t *testing.T) {
	type args struct {
		secretsConfig      *models.SecretsConfig
		gitleaksBinaryPath string
	}
	tests := []struct {
		name string
		args args
		want secrets.Config
	}{
		{
			name: "no config",
			args: args{
				secretsConfig: nil,
			},
			want: secrets.Config{
				Enabled: false,
			},
		},
		{
			name: "no config enabled",
			args: args{
				secretsConfig: &models.SecretsConfig{
					Enabled: nil,
				},
			},
			want: secrets.Config{
				Enabled: false,
			},
		},
		{
			name: "disabled",
			args: args{
				secretsConfig: &models.SecretsConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: secrets.Config{
				Enabled: false,
			},
		},
		{
			name: "enabled",
			args: args{
				secretsConfig: &models.SecretsConfig{
					Enabled: utils.BoolPtr(true),
				},
				gitleaksBinaryPath: "gitleaksBinaryPath",
			},
			want: secrets.Config{
				Enabled:      true,
				ScannersList: []string{"gitleaks"},
				ScannersConfig: &secretscommon.ScannersConfig{
					Gitleaks: gitleaksconfig.Config{
						BinaryPath: "gitleaksBinaryPath",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userSecretsConfigToFamiliesSecretsConfig(tt.args.secretsConfig, tt.args.gitleaksBinaryPath)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("userSecretsConfigToFamiliesSecretsConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_userMalwareConfigToFamiliesMalwareConfig(t *testing.T) {
	type args struct {
		malwareConfig                 *models.MalwareConfig
		clamBinaryPath                string
		freshclamBinaryPath           string
		alternativeFreshclamMirrorURL string
	}
	tests := []struct {
		name string
		args args
		want malware.Config
	}{
		{
			name: "no config",
			args: args{
				malwareConfig: nil,
			},
			want: malware.Config{
				Enabled: false,
			},
		},
		{
			name: "no config enabled",
			args: args{
				malwareConfig: &models.MalwareConfig{
					Enabled: nil,
				},
			},
			want: malware.Config{
				Enabled: false,
			},
		},
		{
			name: "disabled",
			args: args{
				malwareConfig: &models.MalwareConfig{
					Enabled: utils.BoolPtr(false),
				},
			},
			want: malware.Config{
				Enabled: false,
			},
		},
		{
			name: "enabled",
			args: args{
				malwareConfig: &models.MalwareConfig{
					Enabled: utils.BoolPtr(true),
				},
				clamBinaryPath:                "clamscan",
				freshclamBinaryPath:           "freshclam",
				alternativeFreshclamMirrorURL: "",
			},
			want: malware.Config{
				Enabled:      true,
				ScannersList: []string{"clam"},
				ScannersConfig: &malwarecommon.ScannersConfig{
					Clam: malwareconfig.Config{
						ClamScanBinaryPath:            "clamscan",
						FreshclamBinaryPath:           "freshclam",
						AlternativeFreshclamMirrorURL: "",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userMalwareConfigToFamiliesMalwareConfig(
				tt.args.malwareConfig,
				tt.args.clamBinaryPath,
				tt.args.freshclamBinaryPath,
				tt.args.alternativeFreshclamMirrorURL,
			)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("userSecretsConfigToFamiliesSecretsConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
