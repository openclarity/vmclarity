// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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

func (o *Exploit) AsFindingInfo() ExploitFindingInfo {
	return ExploitFindingInfo{
		CveID:       o.CveID,
		Description: o.Description,
		Name:        o.Name,
		ObjectType:  string(ScanFamilyExploit),
		SourceDB:    o.SourceDB,
		Title:       o.Title,
		Urls:        o.Urls,
	}
}

func (o *Exploit) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromExploitFindingInfo(o.AsFindingInfo())
	return *scanFinding
}

func (o *InfoFinder) AsFindingInfo() InfoFinderFindingInfo {
	return InfoFinderFindingInfo{
		Data:       o.Data,
		ObjectType: string(ScanFamilyInfoFinder),
		Path:       o.Path,
		Type:       o.Type,
	}
}

func (o *InfoFinder) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromInfoFinderFindingInfo(o.AsFindingInfo())
	return *scanFinding
}

func (o *Malware) AsFindingInfo() MalwareFindingInfo {
	return MalwareFindingInfo{
		MalwareName: o.MalwareName,
		MalwareType: o.MalwareType,
		ObjectType:  string(ScanFamilyMalware),
		Path:        o.Path,
		RuleName:    o.RuleName,
	}
}

func (o *Malware) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromMalwareFindingInfo(o.AsFindingInfo())
	return *scanFinding
}

func (o *Misconfiguration) AsFindingInfo() MisconfigurationFindingInfo {
	return MisconfigurationFindingInfo{
		Category:    o.Category,
		Description: o.Description,
		Id:          o.Id,
		Location:    o.Location,
		Message:     o.Message,
		ObjectType:  string(ScanFamilyMisconfiguration),
		Remediation: o.Remediation,
		Severity:    o.Severity,
	}
}

func (o *Misconfiguration) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromMisconfigurationFindingInfo(o.AsFindingInfo())
	return *scanFinding
}

func (o *Package) AsFindingInfo() PackageFindingInfo {
	return PackageFindingInfo{
		Cpes:       o.Cpes,
		Language:   o.Language,
		Licenses:   o.Licenses,
		Name:       o.Name,
		ObjectType: string(ScanFamilyPackage),
		Purl:       o.Purl,
		Type:       o.Type,
		Version:    o.Version,
	}
}

func (o *Package) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromPackageFindingInfo(o.AsFindingInfo())
	return *scanFinding
}

func (o *Rootkit) AsFindingInfo() RootkitFindingInfo {
	return RootkitFindingInfo{
		Message:     o.Message,
		ObjectType:  string(ScanFamilyRootkit),
		RootkitName: o.RootkitName,
		RootkitType: o.RootkitType,
	}
}

func (o *Rootkit) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromRootkitFindingInfo(o.AsFindingInfo())
	return *scanFinding
}

func (o *Secret) AsFindingInfo() SecretFindingInfo {
	return SecretFindingInfo{
		Description: o.Description,
		EndColumn:   o.EndColumn,
		EndLine:     o.EndLine,
		FilePath:    o.FilePath,
		Fingerprint: o.Fingerprint,
		ObjectType:  string(ScanFamilySecret),
		StartColumn: o.StartColumn,
		StartLine:   o.StartLine,
	}
}

func (o *Secret) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromSecretFindingInfo(o.AsFindingInfo())
	return *scanFinding
}

func (o *Vulnerability) AsFindingInfo() VulnerabilityFindingInfo {
	return VulnerabilityFindingInfo{
		Cvss:              o.Cvss,
		Description:       o.Description,
		Distro:            o.Distro,
		Fix:               o.Fix,
		LayerId:           o.LayerId,
		Links:             o.Links,
		ObjectType:        string(ScanFamilyVulnerability),
		Package:           o.Package,
		Path:              o.Path,
		Severity:          o.Severity,
		VulnerabilityName: o.VulnerabilityName,
	}
}

func (o *Vulnerability) AsScanFindingInfo() ScanFinding_FindingInfo {
	scanFinding := new(ScanFinding_FindingInfo)
	_ = scanFinding.FromVulnerabilityFindingInfo(o.AsFindingInfo())
	return *scanFinding
}
