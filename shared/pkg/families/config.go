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
	"fmt"
	"net"
	"time"

	"github.com/anchore/syft/syft/source"
	kubeclarityConfig "github.com/openclarity/kubeclarity/shared/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/families/exploits"
	exploitsCommon "github.com/openclarity/vmclarity/shared/pkg/families/exploits/common"
	exploitdbConfig "github.com/openclarity/vmclarity/shared/pkg/families/exploits/exploitdb/config"
	"github.com/openclarity/vmclarity/shared/pkg/families/malware"
	malwareconfig "github.com/openclarity/vmclarity/shared/pkg/families/malware/clam/config"
	malwarecommon "github.com/openclarity/vmclarity/shared/pkg/families/malware/common"
	misconfigurationTypes "github.com/openclarity/vmclarity/shared/pkg/families/misconfiguration/types"
	"github.com/openclarity/vmclarity/shared/pkg/families/rootkits"
	chkrootkitConfig "github.com/openclarity/vmclarity/shared/pkg/families/rootkits/chkrootkit/config"
	rootkitsCommon "github.com/openclarity/vmclarity/shared/pkg/families/rootkits/common"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets/common"
	gitleaksconfig "github.com/openclarity/vmclarity/shared/pkg/families/secrets/gitleaks/config"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
)

const (
	TrivyTimeout       = 300
	GrypeServerTimeout = 2 * time.Minute

	GitleaksBinaryPath            = "GITLEAKS_BINARY_PATH"
	ClamBinaryPath                = "CLAM_BINARY_PATH"
	FreshclamBinaryPath           = "FRESHCLAM_BINARY_PATH"
	AlternativeFreshclamMirrorURL = "ALTERNATIVE_FRESHCLAM_MIRROR_URL"
	LynisInstallPath              = "LYNIS_INSTALL_PATH"
	ChkrootkitBinaryPath          = "CHKROOTKIT_BINARY_PATH"
	ExploitDBAddress              = "EXPLOIT_DB_ADDRESS"
	TrivyServerAddress            = "TRIVY_SERVER_ADDRESS"
	GrypeServerAddress            = "GRYPE_SERVER_ADDRESS"
)

type Config struct {
	// Analyzers
	SBOM sbom.Config `json:"sbom" yaml:"sbom" mapstructure:"sbom"`

	// Scanners
	Vulnerabilities  vulnerabilities.Config       `json:"vulnerabilities" yaml:"vulnerabilities" mapstructure:"vulnerabilities"`
	Secrets          secrets.Config               `json:"secrets" yaml:"secrets" mapstructure:"secrets"`
	Rootkits         rootkits.Config              `json:"rootkits" yaml:"rootkits" mapstructure:"rootkits"`
	Malware          malware.Config               `json:"malware" yaml:"malware" mapstructure:"malware"`
	Misconfiguration misconfigurationTypes.Config `json:"misconfiguration" yaml:"misconfiguration" mapstructure:"misconfiguration"`

	// Enrichers
	Exploits exploits.Config `json:"exploits" yaml:"exploits" mapstructure:"exploits"`
}

type Paths struct {
	// The gitleaks binary path in the scanner image container.
	GitleaksBinaryPath string

	// The clam binary path in the scanner image container.
	ClamBinaryPath string

	// The freshclam binary path in the scanner image container
	FreshclamBinaryPath string

	// The freshclam mirror url to use if it's enabled
	AlternativeFreshclamMirrorURL string

	// The location where Lynis is installed in the scanner image
	LynisInstallPath string

	// The chkrootkit binary path in the scanner image container.
	ChkrootkitBinaryPath string
}

type Addresses struct {
	ExploitsDBAddress  string
	GrypeServerAddress string
	TrivyServerAddress string
}

func setDefaultPaths() {
	// https://github.com/openclarity/vmclarity-tools-base/blob/main/Dockerfile#L33
	viper.SetDefault(GitleaksBinaryPath, "/artifacts/gitleaks")
	// https://github.com/openclarity/vmclarity-tools-base/blob/main/Dockerfile#L35
	viper.SetDefault(LynisInstallPath, "/artifacts/lynis")
	// https://github.com/openclarity/vmclarity-tools-base/blob/main/Dockerfile
	viper.SetDefault(ChkrootkitBinaryPath, "/artifacts/chkrootkit")
	viper.SetDefault(ClamBinaryPath, "clamscan")
	viper.SetDefault(FreshclamBinaryPath, "freshclam")

	viper.AutomaticEnv()
}

func setDefaultAddresses(exploitDBHost string) {
	viper.SetDefault(ExploitDBAddress, fmt.Sprintf("http://%s", net.JoinHostPort(exploitDBHost, "1326")))

	viper.AutomaticEnv()
}

func LoadPaths() Paths {
	setDefaultPaths()
	paths := Paths{
		GitleaksBinaryPath:            viper.GetString(GitleaksBinaryPath),
		ClamBinaryPath:                viper.GetString(ClamBinaryPath),
		FreshclamBinaryPath:           viper.GetString(FreshclamBinaryPath),
		AlternativeFreshclamMirrorURL: viper.GetString(AlternativeFreshclamMirrorURL),
		LynisInstallPath:              viper.GetString(LynisInstallPath),
		ChkrootkitBinaryPath:          viper.GetString(ChkrootkitBinaryPath),
	}

	return paths
}

func LoadAddresses(exploitDBHost string) Addresses {
	setDefaultAddresses(exploitDBHost)
	addresses := Addresses{
		ExploitsDBAddress:  viper.GetString(ExploitDBAddress),
		GrypeServerAddress: viper.GetString(GrypeServerAddress),
		TrivyServerAddress: viper.GetString(TrivyServerAddress),
	}

	return addresses
}

func CreateFamilyConfigFromModel(scanFamiliesConfig *models.ScanFamiliesConfig,
	addresses Addresses,
	familiesPaths Paths,
) Config {
	return Config{
		SBOM:            userSBOMConfigToFamiliesSbomConfig(scanFamiliesConfig.Sbom),
		Vulnerabilities: userVulnConfigToFamiliesVulnConfig(scanFamiliesConfig.Vulnerabilities, addresses.TrivyServerAddress, addresses.GrypeServerAddress),
		Secrets:         userSecretsConfigToFamiliesSecretsConfig(scanFamiliesConfig.Secrets, familiesPaths.GitleaksBinaryPath),
		Exploits:        userExploitsConfigToFamiliesExploitsConfig(scanFamiliesConfig.Exploits, addresses.ExploitsDBAddress),
		Malware: userMalwareConfigToFamiliesMalwareConfig(
			scanFamiliesConfig.Malware,
			familiesPaths.ClamBinaryPath,
			familiesPaths.FreshclamBinaryPath,
			familiesPaths.AlternativeFreshclamMirrorURL,
		),
		Misconfiguration: userMisconfigurationConfigToFamiliesMisconfigurationConfig(scanFamiliesConfig.Misconfigurations, familiesPaths.LynisInstallPath),
		Rootkits:         userRootkitsConfigToFamiliesRootkitsConfig(scanFamiliesConfig.Rootkits, familiesPaths.ChkrootkitBinaryPath),
	}
}

func userRootkitsConfigToFamiliesRootkitsConfig(rootkitsConfig *models.RootkitsConfig, chkRootkitBinaryPath string) rootkits.Config {
	if rootkitsConfig == nil || rootkitsConfig.Enabled == nil || !*rootkitsConfig.Enabled {
		return rootkits.Config{}
	}

	return rootkits.Config{
		Enabled:      true,
		ScannersList: []string{"chkrootkit"},
		Inputs:       nil,
		ScannersConfig: &rootkitsCommon.ScannersConfig{
			Chkrootkit: chkrootkitConfig.Config{
				BinaryPath: chkRootkitBinaryPath,
			},
		},
	}
}

func userSecretsConfigToFamiliesSecretsConfig(secretsConfig *models.SecretsConfig, gitleaksBinaryPath string) secrets.Config {
	if secretsConfig == nil || secretsConfig.Enabled == nil || !*secretsConfig.Enabled {
		return secrets.Config{}
	}
	return secrets.Config{
		Enabled: true,
		// TODO(idanf) This choice should come from the user's configuration
		ScannersList: []string{"gitleaks"},
		Inputs:       nil, // rootfs directory will be determined by the CLI after mount.
		ScannersConfig: &common.ScannersConfig{
			Gitleaks: gitleaksconfig.Config{
				BinaryPath: gitleaksBinaryPath,
			},
		},
	}
}

func userSBOMConfigToFamiliesSbomConfig(sbomConfig *models.SBOMConfig) sbom.Config {
	if sbomConfig == nil || sbomConfig.Enabled == nil || !*sbomConfig.Enabled {
		return sbom.Config{}
	}
	return sbom.Config{
		Enabled: true,
		// TODO(sambetts) This choice should come from the user's configuration
		AnalyzersList: []string{"syft", "trivy"},
		Inputs:        nil, // rootfs directory will be determined by the CLI after mount.
		AnalyzersConfig: &kubeclarityConfig.Config{
			// TODO(sambetts) The user needs to be able to provide this configuration
			Registry: &kubeclarityConfig.Registry{},
			Analyzer: &kubeclarityConfig.Analyzer{
				OutputFormat: "cyclonedx",
				TrivyConfig: kubeclarityConfig.AnalyzerTrivyConfig{
					Timeout: TrivyTimeout,
				},
			},
		},
	}
}

func userMisconfigurationConfigToFamiliesMisconfigurationConfig(misconfigurationConfig *models.MisconfigurationsConfig, lynisInstallPath string) misconfigurationTypes.Config {
	if misconfigurationConfig == nil || misconfigurationConfig.Enabled == nil || !*misconfigurationConfig.Enabled {
		return misconfigurationTypes.Config{}
	}
	return misconfigurationTypes.Config{
		Enabled: true,
		// TODO(sambetts) This choice should come from the user's configuration
		ScannersList: []string{"lynis"},
		Inputs:       nil, // rootfs directory will be determined by the CLI after mount.
		ScannersConfig: misconfigurationTypes.ScannersConfig{
			// TODO(sambetts) Add scanner configurations here as we add them like Lynis
			Lynis: misconfigurationTypes.LynisConfig{
				InstallPath: lynisInstallPath,
			},
		},
	}
}

func userVulnConfigToFamiliesVulnConfig(vulnerabilitiesConfig *models.VulnerabilitiesConfig, trivyServerAddr string, grypeServerAddr string) vulnerabilities.Config {
	if vulnerabilitiesConfig == nil || vulnerabilitiesConfig.Enabled == nil || !*vulnerabilitiesConfig.Enabled {
		return vulnerabilities.Config{}
	}

	var grypeConfig kubeclarityConfig.GrypeConfig
	if grypeServerAddr != "" {
		grypeConfig = kubeclarityConfig.GrypeConfig{
			Mode: kubeclarityConfig.ModeRemote,
			RemoteGrypeConfig: kubeclarityConfig.RemoteGrypeConfig{
				GrypeServerAddress: grypeServerAddr,
				GrypeServerTimeout: GrypeServerTimeout,
			},
		}
	} else {
		grypeConfig = kubeclarityConfig.GrypeConfig{
			Mode: kubeclarityConfig.ModeLocal,
			LocalGrypeConfig: kubeclarityConfig.LocalGrypeConfig{
				UpdateDB:   true,
				DBRootDir:  "/tmp/",
				ListingURL: "https://toolbox-data.anchore.io/grype/databases/listing.json",
				Scope:      source.SquashedScope,
			},
		}
	}

	return vulnerabilities.Config{
		Enabled: true,
		// TODO(sambetts) This choice should come from the user's configuration
		ScannersList:  []string{"grype", "trivy"},
		InputFromSbom: false, // will be determined by the CLI.
		ScannersConfig: &kubeclarityConfig.Config{
			// TODO(sambetts) The user needs to be able to provide this configuration
			Registry: &kubeclarityConfig.Registry{},
			Scanner: &kubeclarityConfig.Scanner{
				GrypeConfig: grypeConfig,
				TrivyConfig: kubeclarityConfig.ScannerTrivyConfig{
					Timeout:    TrivyTimeout,
					ServerAddr: trivyServerAddr,
				},
			},
		},
	}
}

func userExploitsConfigToFamiliesExploitsConfig(exploitsConfig *models.ExploitsConfig, baseURL string) exploits.Config {
	if exploitsConfig == nil || exploitsConfig.Enabled == nil || !*exploitsConfig.Enabled {
		return exploits.Config{}
	}
	// TODO(erezf) Some choices should come from the user's configuration
	return exploits.Config{
		Enabled:       true,
		ScannersList:  []string{"exploitdb"},
		InputFromVuln: true,
		ScannersConfig: &exploitsCommon.ScannersConfig{
			ExploitDB: exploitdbConfig.Config{
				BaseURL: baseURL,
			},
		},
	}
}

func userMalwareConfigToFamiliesMalwareConfig(
	malwareConfig *models.MalwareConfig,
	clamBinaryPath string,
	freshclamBinaryPath string,
	alternativeFreshclamMirrorURL string,
) malware.Config {
	if malwareConfig == nil || malwareConfig.Enabled == nil || !*malwareConfig.Enabled {
		return malware.Config{}
	}

	log.Debugf("clam binary path: %s", clamBinaryPath)
	return malware.Config{
		Enabled:      true,
		ScannersList: []string{"clam"},
		Inputs:       nil, // rootfs directory will be determined by the CLI after mount.
		ScannersConfig: &malwarecommon.ScannersConfig{
			Clam: malwareconfig.Config{
				ClamScanBinaryPath:            clamBinaryPath,
				FreshclamBinaryPath:           freshclamBinaryPath,
				AlternativeFreshclamMirrorURL: alternativeFreshclamMirrorURL,
			},
		},
	}
}
