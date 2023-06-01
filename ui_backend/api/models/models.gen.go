// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package models

import (
	"time"
)

// Defines values for AssetType.
const (
	AWSEC2Instance AssetType = "AWS EC2 Instance"
)

// Defines values for FindingType.
const (
	EXPLOIT          FindingType = "EXPLOIT"
	MALWARE          FindingType = "MALWARE"
	MISCONFIGURATION FindingType = "MISCONFIGURATION"
	PACKAGE          FindingType = "PACKAGE"
	ROOTKIT          FindingType = "ROOTKIT"
	SECRET           FindingType = "SECRET"
	VULNERABILITY    FindingType = "VULNERABILITY"
)

// Defines values for MisconfigurationSeverity.
const (
	MisconfigurationHighSeverity   MisconfigurationSeverity = "MisconfigurationHighSeverity"
	MisconfigurationLowSeverity    MisconfigurationSeverity = "MisconfigurationLowSeverity"
	MisconfigurationMediumSeverity MisconfigurationSeverity = "MisconfigurationMediumSeverity"
)

// Defines values for RootkitType.
const (
	APPLICATION RootkitType = "APPLICATION"
	FIRMWARE    RootkitType = "FIRMWARE"
	KERNEL      RootkitType = "KERNEL"
	MEMORY      RootkitType = "MEMORY"
)

// Defines values for VulnerabilitySeverity.
const (
	CRITICAL   VulnerabilitySeverity = "CRITICAL"
	HIGH       VulnerabilitySeverity = "HIGH"
	LOW        VulnerabilitySeverity = "LOW"
	MEDIUM     VulnerabilitySeverity = "MEDIUM"
	NEGLIGIBLE VulnerabilitySeverity = "NEGLIGIBLE"
)

// ApiResponse An object that is returned in all cases of failures.
type ApiResponse struct {
	Message *string `json:"message,omitempty"`
}

// AssetInfo defines model for AssetInfo.
type AssetInfo struct {
	Location *string    `json:"location,omitempty"`
	Name     *string    `json:"name,omitempty"`
	Type     *AssetType `json:"type,omitempty"`
}

// AssetType defines model for AssetType.
type AssetType string

// Exploit defines model for Exploit.
type Exploit struct {
	CveID       *string   `json:"cveID,omitempty"`
	Description *string   `json:"description,omitempty"`
	Name        *string   `json:"name,omitempty"`
	SourceDB    *string   `json:"sourceDB,omitempty"`
	Title       *string   `json:"title,omitempty"`
	Urls        *[]string `json:"urls"`
}

// ExploitFindingImpact defines model for ExploitFindingImpact.
type ExploitFindingImpact struct {
	AffectedAssetsCount *int     `json:"affectedAssetsCount,omitempty"`
	Exploit             *Exploit `json:"exploit,omitempty"`
}

// FindingTrend Represents the total number of findings at a specific time
type FindingTrend struct {
	Count *int       `json:"count,omitempty"`
	Time  *time.Time `json:"time,omitempty"`
}

// FindingTrends List of the total number of findings per time slot for the specific findingType.
type FindingTrends struct {
	FindingType *FindingType    `json:"findingType,omitempty"`
	Trends      *[]FindingTrend `json:"trends,omitempty"`
}

// FindingType defines model for FindingType.
type FindingType string

// FindingsCount total count of each finding type
type FindingsCount struct {
	Exploits          *int `json:"exploits,omitempty"`
	Malware           *int `json:"malware,omitempty"`
	Misconfigurations *int `json:"misconfigurations,omitempty"`
	Rootkits          *int `json:"rootkits,omitempty"`
	Secrets           *int `json:"secrets,omitempty"`
	Vulnerabilities   *int `json:"vulnerabilities,omitempty"`
}

// FindingsImpact defines model for FindingsImpact.
type FindingsImpact struct {
	// Exploits Top 5 exploit findings sorted by impacted assets count
	Exploits *[]ExploitFindingImpact `json:"exploits,omitempty"`

	// Malware Top 5 malware findings sorted by impacted assets count
	Malware *[]MalwareFindingImpact `json:"malware,omitempty"`

	// Misconfigurations Top 5 misconfiguration findings sorted by impacted assets count
	Misconfigurations *[]MisconfigurationFindingImpact `json:"misconfigurations,omitempty"`

	// Packages Top 5 package findings sorted by impacted assets count
	Packages *[]PackageFindingImpact `json:"packages,omitempty"`

	// Rootkits Top 5 rootkit findings sorted by impacted assets count
	Rootkits *[]RootkitFindingImpact `json:"rootkits,omitempty"`

	// Secrets Top 5 secret findings sorted by impacted assets count
	Secrets *[]SecretFindingImpact `json:"secrets,omitempty"`

	// Vulnerabilities Top 5 vulnerability findings sorted by impacted assets count
	Vulnerabilities *[]VulnerabilityFindingImpact `json:"vulnerabilities,omitempty"`
}

// FindingsTrends List of finding trends for all finding types.
type FindingsTrends = []FindingTrends

// Malware defines model for Malware.
type Malware struct {
	MalwareName *string      `json:"malwareName,omitempty"`
	MalwareType *MalwareType `json:"malwareType,omitempty"`

	// Path Path of the file that contains malware
	Path *string `json:"path,omitempty"`
}

// MalwareFindingImpact defines model for MalwareFindingImpact.
type MalwareFindingImpact struct {
	AffectedAssetsCount *int     `json:"affectedAssetsCount,omitempty"`
	Malware             *Malware `json:"malware,omitempty"`
}

// MalwareType defines model for MalwareType.
type MalwareType = string

// Misconfiguration defines model for Misconfiguration.
type Misconfiguration struct {
	Message         *string                   `json:"message,omitempty"`
	Remediation     *string                   `json:"remediation,omitempty"`
	ScannedPath     *string                   `json:"scannedPath,omitempty"`
	ScannerName     *string                   `json:"scannerName,omitempty"`
	Severity        *MisconfigurationSeverity `json:"severity,omitempty"`
	TestCategory    *string                   `json:"testCategory,omitempty"`
	TestDescription *string                   `json:"testDescription,omitempty"`
	TestID          *string                   `json:"testID,omitempty"`
}

// MisconfigurationFindingImpact defines model for MisconfigurationFindingImpact.
type MisconfigurationFindingImpact struct {
	AffectedAssetsCount *int              `json:"affectedAssetsCount,omitempty"`
	Misconfiguration    *Misconfiguration `json:"misconfiguration,omitempty"`
}

// MisconfigurationSeverity defines model for MisconfigurationSeverity.
type MisconfigurationSeverity string

// Package defines model for Package.
type Package struct {
	Name    *string `json:"name,omitempty"`
	Purl    *string `json:"purl,omitempty"`
	Version *string `json:"version,omitempty"`
}

// PackageFindingImpact defines model for PackageFindingImpact.
type PackageFindingImpact struct {
	AffectedAssetsCount *int     `json:"affectedAssetsCount,omitempty"`
	Package             *Package `json:"package,omitempty"`
}

// RegionFindings Total findings for a region
type RegionFindings struct {
	// FindingsCount total count of each finding type
	FindingsCount *FindingsCount `json:"findingsCount,omitempty"`
	RegionName    *string        `json:"regionName,omitempty"`
}

// RiskiestAssets defines model for RiskiestAssets.
type RiskiestAssets struct {
	// Exploits Top 5 riskiest assets sorted by exploits count
	Exploits *[]RiskyAsset `json:"exploits,omitempty"`

	// Malware Top 5 riskiest assets sorted by malware count
	Malware *[]RiskyAsset `json:"malware,omitempty"`

	// Misconfigurations Top 5 riskiest assets sorted by misconfigurations count
	Misconfigurations *[]RiskyAsset `json:"misconfigurations,omitempty"`

	// Rootkits Top 5 riskiest assets sorted by rootkits count
	Rootkits *[]RiskyAsset `json:"rootkits,omitempty"`

	// Secrets Top 5 riskiest assets sorted by secrets count
	Secrets *[]RiskyAsset `json:"secrets,omitempty"`

	// Vulnerabilities Top 5 riskiest assets sorted by vulnerabilities
	Vulnerabilities *[]VulnerabilityRiskyAsset `json:"vulnerabilities,omitempty"`
}

// RiskiestRegions defines model for RiskiestRegions.
type RiskiestRegions struct {
	// Regions List of regions with the findings that was found on them. Regions with no findings will not be reported.
	Regions *[]RegionFindings `json:"regions,omitempty"`
}

// RiskyAsset Total number of findings for an asset
type RiskyAsset struct {
	AssetInfo *AssetInfo `json:"assetInfo,omitempty"`
	Count     *int       `json:"count,omitempty"`
}

// Rootkit defines model for Rootkit.
type Rootkit struct {
	RootkitName *string      `json:"rootkitName,omitempty"`
	RootkitType *RootkitType `json:"rootkitType,omitempty"`
}

// RootkitFindingImpact defines model for RootkitFindingImpact.
type RootkitFindingImpact struct {
	AffectedAssetsCount *int     `json:"affectedAssetsCount,omitempty"`
	Rootkit             *Rootkit `json:"rootkit,omitempty"`
}

// RootkitType defines model for RootkitType.
type RootkitType string

// Secret defines model for Secret.
type Secret struct {
	EndColumn *int `json:"endColumn,omitempty"`
	EndLine   *int `json:"endLine,omitempty"`

	// FilePath Name of the file containing the secret
	FilePath *string `json:"filePath,omitempty"`

	// Fingerprint Note: this is not unique
	Fingerprint *string `json:"fingerprint,omitempty"`
	StartColumn *int    `json:"startColumn,omitempty"`
	StartLine   *int    `json:"startLine,omitempty"`
}

// SecretFindingImpact defines model for SecretFindingImpact.
type SecretFindingImpact struct {
	AffectedAssetsCount *int    `json:"affectedAssetsCount,omitempty"`
	Secret              *Secret `json:"secret,omitempty"`
}

// VulnerabilitiesFindingImpact defines model for VulnerabilitiesFindingImpact.
type VulnerabilitiesFindingImpact = []VulnerabilityFindingImpact

// Vulnerability defines model for Vulnerability.
type Vulnerability struct {
	Cvss              *[]VulnerabilityCvss   `json:"cvss"`
	Severity          *VulnerabilitySeverity `json:"severity,omitempty"`
	VulnerabilityName *string                `json:"vulnerabilityName,omitempty"`
}

// VulnerabilityCvss defines model for VulnerabilityCvss.
type VulnerabilityCvss struct {
	Metrics *VulnerabilityCvssMetrics `json:"metrics,omitempty"`
	Vector  *string                   `json:"vector,omitempty"`
	Version *string                   `json:"version,omitempty"`
}

// VulnerabilityCvssMetrics defines model for VulnerabilityCvssMetrics.
type VulnerabilityCvssMetrics struct {
	BaseScore           *float32 `json:"baseScore,omitempty"`
	ExploitabilityScore *float32 `json:"exploitabilityScore,omitempty"`
	ImpactScore         *float32 `json:"impactScore,omitempty"`
}

// VulnerabilityFindingImpact defines model for VulnerabilityFindingImpact.
type VulnerabilityFindingImpact struct {
	AffectedAssetsCount *int           `json:"affectedAssetsCount,omitempty"`
	Vulnerability       *Vulnerability `json:"vulnerability,omitempty"`
}

// VulnerabilityRiskyAsset Total number of vulnerability findings for an asset
type VulnerabilityRiskyAsset struct {
	AssetInfo                      *AssetInfo `json:"assetInfo,omitempty"`
	CriticalVulnerabilitiesCount   *int       `json:"criticalVulnerabilitiesCount,omitempty"`
	HighVulnerabilitiesCount       *int       `json:"highVulnerabilitiesCount,omitempty"`
	LowVulnerabilitiesCount        *int       `json:"lowVulnerabilitiesCount,omitempty"`
	MediumVulnerabilitiesCount     *int       `json:"mediumVulnerabilitiesCount,omitempty"`
	NegligibleVulnerabilitiesCount *int       `json:"negligibleVulnerabilitiesCount,omitempty"`
}

// VulnerabilitySeverity defines model for VulnerabilitySeverity.
type VulnerabilitySeverity string

// EndTime defines model for endTime.
type EndTime = time.Time

// ExampleFilter defines model for exampleFilter.
type ExampleFilter = string

// StartTime defines model for startTime.
type StartTime = time.Time

// UnknownError An object that is returned in all cases of failures.
type UnknownError = ApiResponse

// GetDashboardFindingsTrendsParams defines parameters for GetDashboardFindingsTrends.
type GetDashboardFindingsTrendsParams struct {
	StartTime StartTime `form:"startTime" json:"startTime"`
	EndTime   EndTime   `form:"endTime" json:"endTime"`
}
