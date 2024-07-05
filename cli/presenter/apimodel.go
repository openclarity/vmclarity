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

package presenter

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
	exploits "github.com/openclarity/vmclarity/scanner/families/exploits/types"
	infofinder "github.com/openclarity/vmclarity/scanner/families/infofinder/types"
	malware "github.com/openclarity/vmclarity/scanner/families/malware/types"
	misconfiguration "github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	rootkits "github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	rootkitsTypes "github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	sbom "github.com/openclarity/vmclarity/scanner/families/sbom/types"
	secrets "github.com/openclarity/vmclarity/scanner/families/secrets/types"
	vulnerabilities "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/types"
	"github.com/openclarity/vmclarity/scanner/utils/cyclonedx_helper"
	"github.com/openclarity/vmclarity/scanner/utils/vulnerability"
)

func ConvertSBOMResultToPackages(sbomResults *sbom.Result) []apitypes.Package {
	packages := []apitypes.Package{}

	if sbomResults == nil || sbomResults.SBOM == nil || sbomResults.SBOM.Components == nil {
		return packages
	}

	for _, component := range *sbomResults.SBOM.Components {
		packages = append(packages, apitypes.Package{
			Cpes:     to.Ptr([]string{component.CPE}),
			Language: to.Ptr(cyclonedx_helper.GetComponentLanguage(component)),
			Licenses: to.Ptr(cyclonedx_helper.GetComponentLicenses(component)),
			Name:     to.Ptr(component.Name),
			Purl:     to.Ptr(component.PackageURL),
			Type:     to.Ptr(string(component.Type)),
			Version:  to.Ptr(component.Version),
		})
	}

	return packages
}

func ConvertVulnResultToVulnerabilities(vulnerabilitiesResults *vulnerabilities.Result) []apitypes.Vulnerability {
	vuls := []apitypes.Vulnerability{}

	if vulnerabilitiesResults == nil || vulnerabilitiesResults.MergedVulnerabilitiesByKey == nil {
		return vuls
	}

	for _, vulCandidates := range vulnerabilitiesResults.MergedVulnerabilitiesByKey {
		if len(vulCandidates) < 1 {
			continue
		}

		vulCandidate := vulCandidates[0]

		vul := apitypes.Vulnerability{
			Cvss:              ConvertVulnCvssToAPIModel(vulCandidate.Vulnerability.CVSS),
			Description:       to.Ptr(vulCandidate.Vulnerability.Description),
			Distro:            ConvertVulnDistroToAPIModel(vulCandidate.Vulnerability.Distro),
			Fix:               ConvertVulnFixToAPIModel(vulCandidate.Vulnerability.Fix),
			LayerId:           to.Ptr(vulCandidate.Vulnerability.LayerID),
			Links:             to.Ptr(vulCandidate.Vulnerability.Links),
			Package:           ConvertVulnPackageToAPIModel(vulCandidate.Vulnerability.Package),
			Path:              to.Ptr(vulCandidate.Vulnerability.Path),
			Severity:          ConvertVulnSeverityToAPIModel(vulCandidate.Vulnerability.Severity),
			VulnerabilityName: to.Ptr(vulCandidate.Vulnerability.ID),
		}
		vuls = append(vuls, vul)
	}

	return vuls
}

func ConvertVulnSeverityToAPIModel(severity string) *apitypes.VulnerabilitySeverity {
	switch strings.ToUpper(severity) {
	case vulnerability.DEFCON1, vulnerability.CRITICAL:
		return to.Ptr(apitypes.CRITICAL)
	case vulnerability.HIGH:
		return to.Ptr(apitypes.HIGH)
	case vulnerability.MEDIUM:
		return to.Ptr(apitypes.MEDIUM)
	case vulnerability.LOW:
		return to.Ptr(apitypes.LOW)
	case vulnerability.NEGLIGIBLE, vulnerability.UNKNOWN, vulnerability.NONE:
		return to.Ptr(apitypes.NEGLIGIBLE)
	default:
		log.Errorf("Can't convert severity %q, treating as negligible", severity)
		return to.Ptr(apitypes.NEGLIGIBLE)
	}
}

func ConvertVulnFixToAPIModel(fix vulnerabilities.Fix) *apitypes.VulnerabilityFix {
	return &apitypes.VulnerabilityFix{
		State:    to.Ptr(fix.State),
		Versions: to.Ptr(fix.Versions),
	}
}

func ConvertVulnDistroToAPIModel(distro vulnerabilities.Distro) *apitypes.VulnerabilityDistro {
	return &apitypes.VulnerabilityDistro{
		IDLike:  to.Ptr(distro.IDLike),
		Name:    to.Ptr(distro.Name),
		Version: to.Ptr(distro.Version),
	}
}

func ConvertVulnPackageToAPIModel(p vulnerabilities.Package) *apitypes.Package {
	return &apitypes.Package{
		Cpes:     to.Ptr(p.CPEs),
		Language: to.Ptr(p.Language),
		Licenses: to.Ptr(p.Licenses),
		Name:     to.Ptr(p.Name),
		Purl:     to.Ptr(p.PURL),
		Type:     to.Ptr(p.Type),
		Version:  to.Ptr(p.Version),
	}
}

func ConvertVulnCvssToAPIModel(cvss []vulnerabilities.CVSS) *[]apitypes.VulnerabilityCvss {
	if cvss == nil {
		return nil
	}
	// nolint:prealloc
	var ret []apitypes.VulnerabilityCvss
	for _, c := range cvss {
		var exploitabilityScore *float32
		if c.Metrics.ExploitabilityScore != nil {
			exploitabilityScore = to.Ptr[float32](float32(*c.Metrics.ExploitabilityScore))
		}
		var impactScore *float32
		if c.Metrics.ImpactScore != nil {
			impactScore = to.Ptr[float32](float32(*c.Metrics.ImpactScore))
		}
		ret = append(ret, apitypes.VulnerabilityCvss{
			Metrics: &apitypes.VulnerabilityCvssMetrics{
				BaseScore:           to.Ptr(float32(c.Metrics.BaseScore)),
				ExploitabilityScore: exploitabilityScore,
				ImpactScore:         impactScore,
			},
			Vector:  to.Ptr(c.Vector),
			Version: to.Ptr(c.Version),
		})
	}

	return &ret
}

func ConvertMalwareResultToMalwareAndMetadata(malwareResults *malware.Result) ([]apitypes.Malware, []apitypes.ScannerMetadata) {
	malwareList := []apitypes.Malware{}
	metadata := []apitypes.ScannerMetadata{}

	if malwareResults == nil || malwareResults.Malwares == nil {
		return malwareList, metadata
	}

	for _, m := range malwareResults.Malwares {
		mal := m // Prevent loop variable pointer export
		malwareList = append(malwareList, apitypes.Malware{
			MalwareName: to.PtrOrNil(mal.MalwareName),
			MalwareType: to.PtrOrNil(mal.MalwareType),
			Path:        to.PtrOrNil(mal.Path),
			RuleName:    to.PtrOrNil(mal.RuleName),
		})
	}

	for n, s := range malwareResults.ScansSummary {
		name, summary := n, s // Prevent loop variable pointer export
		metadata = append(metadata, apitypes.ScannerMetadata{
			ScannerName: &name,
			ScannerSummary: &apitypes.ScannerSummary{
				DataRead:           &summary.DataRead,
				DataScanned:        &summary.DataScanned,
				EngineVersion:      &summary.EngineVersion,
				InfectedFiles:      &summary.InfectedFiles,
				KnownViruses:       &summary.KnownViruses,
				ScannedDirectories: &summary.ScannedDirectories,
				ScannedFiles:       &summary.ScannedFiles,
				SuspectedFiles:     &summary.SuspectedFiles,
				TimeTaken:          &summary.TimeTaken,
			},
		})
	}

	return malwareList, metadata
}

func ConvertSecretsResultToSecrets(secretsResults *secrets.Result) []apitypes.Secret {
	secretsSlice := []apitypes.Secret{}

	if secretsResults == nil || secretsResults.Findings == nil {
		return secretsSlice
	}

	for i := range secretsResults.Findings {
		finding := secretsResults.Findings[i]
		secretsSlice = append(secretsSlice, apitypes.Secret{
			Description: &finding.Description,
			EndLine:     &finding.EndLine,
			FilePath:    &finding.File,
			Fingerprint: &finding.Fingerprint,
			StartLine:   &finding.StartLine,
			StartColumn: &finding.StartColumn,
			EndColumn:   &finding.EndColumn,
		})
	}

	return secretsSlice
}

func ConvertExploitsResultToExploits(exploitsResults *exploits.Exploits) []apitypes.Exploit {
	retExploits := []apitypes.Exploit{}

	if exploitsResults == nil || exploitsResults.MergedExploits == nil {
		return retExploits
	}

	for i := range exploitsResults.MergedExploits {
		exploit := exploitsResults.MergedExploits[i]
		retExploits = append(retExploits, apitypes.Exploit{
			CveID:       &exploit.CveID,
			Description: &exploit.Description,
			Name:        &exploit.Name,
			SourceDB:    &exploit.SourceDB,
			Title:       &exploit.Title,
			Urls:        &exploit.URLs,
		})
	}

	return retExploits
}

func MisconfigurationSeverityToAPIMisconfigurationSeverity(sev misconfiguration.Severity) (apitypes.MisconfigurationSeverity, error) {
	switch sev {
	case misconfiguration.HighSeverity:
		return apitypes.MisconfigurationHighSeverity, nil
	case misconfiguration.MediumSeverity:
		return apitypes.MisconfigurationMediumSeverity, nil
	case misconfiguration.LowSeverity:
		return apitypes.MisconfigurationLowSeverity, nil
	case misconfiguration.InfoSeverity:
		return apitypes.MisconfigurationInfoSeverity, nil
	default:
		return apitypes.MisconfigurationLowSeverity, fmt.Errorf("unknown severity level %v", sev)
	}
}

func ConvertMisconfigurationResultToMisconfigurationsAndScanners(misconfigurationResults *misconfiguration.Result) ([]apitypes.Misconfiguration, []string, error) {
	if misconfigurationResults == nil || misconfigurationResults.Misconfigurations == nil {
		return nil, nil, nil
	}

	misconfigurations := []apitypes.Misconfiguration{}
	scanners := misconfigurationResults.Metadata.Scanners()

	for i := range misconfigurationResults.Misconfigurations {
		// create a separate variable for the loop because we need
		// pointers for the API model and we can't safely take pointers
		// to a loop variable.
		misconfig := misconfigurationResults.Misconfigurations[i]

		severity, err := MisconfigurationSeverityToAPIMisconfigurationSeverity(misconfig.Severity)
		if err != nil {
			return misconfigurations, scanners, fmt.Errorf("unable to convert scanner result severity to API severity: %w", err)
		}

		misconfigurations = append(misconfigurations, apitypes.Misconfiguration{
			ScannerName: &misconfig.ScannerName,
			Location:    &misconfig.Location,
			Category:    &misconfig.Category,
			Id:          &misconfig.ID,
			Description: &misconfig.Description,
			Severity:    &severity,
			Message:     &misconfig.Message,
			Remediation: &misconfig.Remediation,
		})
	}

	return misconfigurations, scanners, nil
}

func ConvertInfoFinderResultToInfosAndScanners(results *infofinder.Infos) ([]apitypes.InfoFinderInfo, []string, error) {
	infos := []apitypes.InfoFinderInfo{}

	if results == nil || results.FlattenedInfos == nil {
		return infos, []string{}, nil
	}

	for i := range results.FlattenedInfos {
		info := results.FlattenedInfos[i]

		infos = append(infos, apitypes.InfoFinderInfo{
			Data:        &info.Data,
			Path:        &info.Path,
			ScannerName: &info.ScannerName,
			Type:        convertInfoTypeToAPIModel(info.Type),
		})
	}

	return infos, results.Metadata.Scanners(), nil
}

func convertInfoTypeToAPIModel(infoType infofinder.InfoType) *apitypes.InfoType {
	switch infoType {
	case infofinder.SSHKnownHostFingerprint:
		return to.Ptr(apitypes.InfoTypeSSHKnownHostFingerprint)
	case infofinder.SSHAuthorizedKeyFingerprint:
		return to.Ptr(apitypes.InfoTypeSSHAuthorizedKeyFingerprint)
	case infofinder.SSHPrivateKeyFingerprint:
		return to.Ptr(apitypes.InfoTypeSSHPrivateKeyFingerprint)
	case infofinder.SSHDaemonKeyFingerprint:
		return to.Ptr(apitypes.InfoTypeSSHDaemonKeyFingerprint)
	default:
		log.Errorf("Can't convert info type %q, treating as %v", infoType, apitypes.InfoTypeUNKNOWN)
		return to.Ptr(apitypes.InfoTypeUNKNOWN)
	}
}

func ConvertRootkitsResultToRootkits(rootkitsResults *rootkits.Result) []apitypes.Rootkit {
	rootkitsList := []apitypes.Rootkit{}

	if rootkitsResults == nil || rootkitsResults.Rootkits == nil {
		return rootkitsList
	}

	for i := range rootkitsResults.Rootkits {
		rootkit := rootkitsResults.Rootkits[i]
		rootkitsList = append(rootkitsList, apitypes.Rootkit{
			Message:     &rootkit.Message,
			RootkitName: &rootkit.RootkitName,
			RootkitType: ConvertRootkitTypeToAPIModel(rootkit.RootkitType),
		})
	}

	return rootkitsList
}

func ConvertRootkitTypeToAPIModel(rootkitType rootkitsTypes.RootkitType) *apitypes.RootkitType {
	switch rootkitType {
	case rootkitsTypes.APPLICATION:
		return to.Ptr(apitypes.RootkitTypeAPPLICATION)
	case rootkitsTypes.FIRMWARE:
		return to.Ptr(apitypes.RootkitTypeFIRMWARE)
	case rootkitsTypes.KERNEL:
		return to.Ptr(apitypes.RootkitTypeKERNEL)
	case rootkitsTypes.MEMORY:
		return to.Ptr(apitypes.RootkitTypeMEMORY)
	case rootkitsTypes.UNKNOWN:
		return to.Ptr(apitypes.RootkitTypeUNKNOWN)
	default:
		log.Errorf("Can't convert rootkit type %q, treating as %v", rootkitType, apitypes.RootkitTypeUNKNOWN)
		return to.Ptr(apitypes.RootkitTypeUNKNOWN)
	}
}
