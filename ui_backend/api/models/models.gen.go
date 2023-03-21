// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.3 DO NOT EDIT.
package models

// Defines values for AssetType.
const (
	AWSEC2Instance AssetType = "AWS EC2 Instance"
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

// FindingsCount total count of each finding type
type FindingsCount struct {
	Exploits          *int `json:"exploits,omitempty"`
	Malware           *int `json:"malware,omitempty"`
	Misconfigurations *int `json:"misconfigurations,omitempty"`
	Rootkits          *int `json:"rootkits,omitempty"`
	Secrets           *int `json:"secrets,omitempty"`
	Vulnerabilities   *int `json:"vulnerabilities,omitempty"`
}

// RegionFindings Total findings for a region
type RegionFindings struct {
	// FindingsCount total count of each finding type
	FindingsCount *FindingsCount `json:"findingsCount,omitempty"`
	RegionName    *string        `json:"regionName,omitempty"`
}

// RiskiestAssets defines model for RiskiestAssets.
type RiskiestAssets struct {
	// Exploits Top 5 assets sorted by max exploits count
	Exploits *[]RiskyAsset `json:"exploits,omitempty"`

	// Malware Top 5 assets sorted by max malware count
	Malware *[]RiskyAsset `json:"malware,omitempty"`

	// Misconfigurations Top 5 assets sorted by max misconfigurations count
	Misconfigurations *[]RiskyAsset `json:"misconfigurations,omitempty"`

	// Rootkits Top 5 assets sorted by max rootkits count
	Rootkits *[]RiskyAsset `json:"rootkits,omitempty"`

	// Secrets Top 5 assets sorted by max secrets count
	Secrets *[]RiskyAsset `json:"secrets,omitempty"`

	// Vulnerabilities Top 5 assets sorted by max vulnerabilites count
	Vulnerabilities *[]VulnerabilityRiskyAsset `json:"vulnerabilities,omitempty"`
}

// RiskiestRegions defines model for RiskiestRegions.
type RiskiestRegions struct {
	// Regions List of regions with the findings that was found on them. Regions with no findings will not be reported.
	Regions *[]RegionFindings `json:"regions,omitempty"`
}

// RiskyAsset Total findings for an asset
type RiskyAsset struct {
	AssetInfo *AssetInfo `json:"assetInfo,omitempty"`
	Count     *int       `json:"count,omitempty"`
}

// VulnerabilityCount defines model for VulnerabilityCount.
type VulnerabilityCount struct {
	CriticalVulnerabilitiesCount   *int `json:"criticalVulnerabilitiesCount,omitempty"`
	HighVulnerabilitiesCount       *int `json:"highVulnerabilitiesCount,omitempty"`
	LowVulnerabilitiesCount        *int `json:"lowVulnerabilitiesCount,omitempty"`
	MediumVulnerabilitiesCount     *int `json:"mediumVulnerabilitiesCount,omitempty"`
	NegligibleVulnerabilitiesCount *int `json:"negligibleVulnerabilitiesCount,omitempty"`
}

// VulnerabilityRiskyAsset Total vulnerabilities findings for an asset
type VulnerabilityRiskyAsset struct {
	AssetInfo *AssetInfo          `json:"assetInfo,omitempty"`
	Count     *VulnerabilityCount `json:"count,omitempty"`
}

// ExampleFilter defines model for exampleFilter.
type ExampleFilter = string

// UnknownError An object that is returned in all cases of failures.
type UnknownError = ApiResponse
