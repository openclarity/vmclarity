// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.3 DO NOT EDIT.
package models

// Defines values for CloudProvider.
const (
	AWS CloudProvider = "AWS"
)

// Defines values for MalwareType.
const (
	ADWARE     MalwareType = "ADWARE"
	RANSOMWARE MalwareType = "RANSOMWARE"
	SPYWARE    MalwareType = "SPYWARE"
	TROJAN     MalwareType = "TROJAN"
	VIRUS      MalwareType = "VIRUS"
	WORM       MalwareType = "WORM"
)

// Defines values for RootkitType.
const (
	APPLICATION RootkitType = "APPLICATION"
	FIRMWARE    RootkitType = "FIRMWARE"
	KERNEL      RootkitType = "KERNEL"
	MEMORY      RootkitType = "MEMORY"
)

// ApiResponse An object that is returned in all cases of failures.
type ApiResponse struct {
	Message *string `json:"message,omitempty"`
}

// CloudProvider defines model for CloudProvider.
type CloudProvider string

// ExploitInfo defines model for ExploitInfo.
type ExploitInfo struct {
	Description     *string   `json:"description,omitempty"`
	Id              *string   `json:"id,omitempty"`
	Vulnerabilities *[]string `json:"vulnerabilities,omitempty"`
}

// ExploitScan defines model for ExploitScan.
type ExploitScan struct {
	Packages *[]ExploitInfo `json:"packages,omitempty"`
}

// Instance defines model for Instance.
type Instance struct {
	Id          *string       `json:"id,omitempty"`
	Instance    *InstanceInfo `json:"instance,omitempty"`
	ScanResults *uint32       `json:"scanResults,omitempty"`
}

// InstanceInfo defines model for InstanceInfo.
type InstanceInfo struct {
	Id               *string        `json:"id,omitempty"`
	InstanceName     *string        `json:"instanceName,omitempty"`
	InstanceProvider *CloudProvider `json:"instanceProvider,omitempty"`
	Location         *string        `json:"location,omitempty"`
}

// MalwareInfo defines model for MalwareInfo.
type MalwareInfo struct {
	Id          *string      `json:"id,omitempty"`
	MalwareName *string      `json:"malwareName,omitempty"`
	MalwareType *MalwareType `json:"malwareType,omitempty"`

	// Path Path of the file that contains malware
	Path *string `json:"path,omitempty"`
}

// MalwareScan defines model for MalwareScan.
type MalwareScan struct {
	Packages *[]MalwareInfo `json:"packages,omitempty"`
}

// MalwareType defines model for MalwareType.
type MalwareType string

// MisconfigurationInfo defines model for MisconfigurationInfo.
type MisconfigurationInfo struct {
	Description *string `json:"description,omitempty"`
	Id          *string `json:"id,omitempty"`

	// Path Path of the file that contains misconfigurations
	Path *string `json:"path,omitempty"`
}

// MisconfigurationScan defines model for MisconfigurationScan.
type MisconfigurationScan struct {
	Packages *[]MisconfigurationInfo `json:"packages,omitempty"`
}

// Package defines model for Package.
type Package struct {
	Id          *string      `json:"id,omitempty"`
	PackageInfo *PackageInfo `json:"packageInfo,omitempty"`
}

// PackageInfo defines model for PackageInfo.
type PackageInfo struct {
	Id             *string `json:"id,omitempty"`
	PackageName    *string `json:"packageName,omitempty"`
	PackageVersion *string `json:"packageVersion,omitempty"`
}

// RootkitInfo defines model for RootkitInfo.
type RootkitInfo struct {
	Id *string `json:"id,omitempty"`

	// Path Path of the file that contains rootkit
	Path        *string      `json:"path,omitempty"`
	RootkitName *string      `json:"rootkitName,omitempty"`
	RootkitType *RootkitType `json:"rootkitType,omitempty"`
}

// RootkitScan defines model for RootkitScan.
type RootkitScan struct {
	Packages *[]RootkitInfo `json:"packages,omitempty"`
}

// RootkitType defines model for RootkitType.
type RootkitType string

// SbomScan defines model for SbomScan.
type SbomScan struct {
	Packages *[]Package `json:"packages,omitempty"`
}

// ScanResults defines model for ScanResults.
type ScanResults struct {
	Exploits          *ExploitScan          `json:"exploits,omitempty"`
	Id                *string               `json:"id,omitempty"`
	Malwares          *VulnerabilityScan    `json:"malwares,omitempty"`
	Misconfigurations *MisconfigurationScan `json:"misconfigurations,omitempty"`
	Sboms             *SbomScan             `json:"sboms,omitempty"`
	Secrets           *SecretScan           `json:"secrets,omitempty"`
	Vulnerabilities   *VulnerabilityScan    `json:"vulnerabilities,omitempty"`
}

// ScanResultsSummary defines model for ScanResultsSummary.
type ScanResultsSummary struct {
	ExploitsCount          *int `json:"exploitsCount,omitempty"`
	MalwaresCount          *int `json:"malwaresCount,omitempty"`
	MisconfigurationsCount *int `json:"misconfigurationsCount,omitempty"`
	PackagesCount          *int `json:"packagesCount,omitempty"`
	SecretsCount           *int `json:"secretsCount,omitempty"`
	VulnerabilitiesCount   *int `json:"vulnerabilitiesCount,omitempty"`
}

// SecretInfo defines model for SecretInfo.
type SecretInfo struct {
	Description *string `json:"description,omitempty"`
	Id          *string `json:"id,omitempty"`

	// Path Path of the file that contains secrets
	Path *string `json:"path,omitempty"`
}

// SecretScan defines model for SecretScan.
type SecretScan struct {
	Packages *[]SecretInfo `json:"packages,omitempty"`
}

// SuccessResponse An object that is returned in cases of success that returns nothing.
type SuccessResponse struct {
	Message *string `json:"message,omitempty"`
}

// Vulnerability defines model for Vulnerability.
type Vulnerability struct {
	Id          *string            `json:"id,omitempty"`
	PackageInfo *VulnerabilityInfo `json:"packageInfo,omitempty"`
}

// VulnerabilityInfo defines model for VulnerabilityInfo.
type VulnerabilityInfo struct {
	Description       *string `json:"description,omitempty"`
	Id                *string `json:"id,omitempty"`
	VulnerabilityName *string `json:"vulnerabilityName,omitempty"`
}

// VulnerabilityScan defines model for VulnerabilityScan.
type VulnerabilityScan struct {
	Vulnerabilities *[]Vulnerability `json:"vulnerabilities,omitempty"`
}

// InstanceID defines model for instanceID.
type InstanceID = string

// Page defines model for page.
type Page = int

// PageSize defines model for pageSize.
type PageSize = int

// ScanID defines model for scanID.
type ScanID = string

// Success An object that is returned in cases of success that returns nothing.
type Success = SuccessResponse

// UnknownError An object that is returned in all cases of failures.
type UnknownError = ApiResponse

// GetInstancesParams defines parameters for GetInstances.
type GetInstancesParams struct {
	// Page Page number of the query
	Page int `form:"page" json:"page"`

	// PageSize Maximum items to return
	PageSize int `form:"pageSize" json:"pageSize"`

	// SortKey Sort key
	SortKey string `form:"sortKey" json:"sortKey"`
}

// GetInstancesInstanceIDScanresultsParams defines parameters for GetInstancesInstanceIDScanresults.
type GetInstancesInstanceIDScanresultsParams struct {
	// Page Page number of the query
	Page int `form:"page" json:"page"`

	// PageSize Maximum items to return
	PageSize int `form:"pageSize" json:"pageSize"`

	// SortKey Sort key
	SortKey string `form:"sortKey" json:"sortKey"`
}

// PostInstancesJSONRequestBody defines body for PostInstances for application/json ContentType.
type PostInstancesJSONRequestBody = InstanceInfo

// PutInstancesInstanceIDJSONRequestBody defines body for PutInstancesInstanceID for application/json ContentType.
type PutInstancesInstanceIDJSONRequestBody = InstanceInfo

// PostInstancesInstanceIDScanresultsJSONRequestBody defines body for PostInstancesInstanceIDScanresults for application/json ContentType.
type PostInstancesInstanceIDScanresultsJSONRequestBody = ScanResults

// PutInstancesInstanceIDScanresultsScanIDJSONRequestBody defines body for PutInstancesInstanceIDScanresultsScanID for application/json ContentType.
type PutInstancesInstanceIDScanresultsScanIDJSONRequestBody = ScanResults
