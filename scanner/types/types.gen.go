// Package types provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package types

import (
	"time"
)

// Defines values for ConfigOutputFormat.
const (
	VMClarityJSON ConfigOutputFormat = "vmclarity-json"
)

// Defines values for StatusState.
const (
	Done     StatusState = "Done"
	Failed   StatusState = "Failed"
	NotReady StatusState = "NotReady"
	Ready    StatusState = "Ready"
	Running  StatusState = "Running"
)

// Config Describes a scanner config.
type Config struct {
	// File The file with the configuration required by the scanner plugin. This is a path on the filesystem to the config file.
	File *string `json:"file,omitempty"`

	// InputDir The directory which should be scanned by the scanner plugin.
	InputDir string `json:"inputDir" validate:"required"`

	// OutputDir The directory where the scanner plugin should store it's findings.
	OutputDir string `json:"outputDir" validate:"required"`

	// OutputFormat The format in which the scanner plugin should store it's findings.
	// To ensure operability with VMClarity API, the format must be one of enum values.
	// However, the scanner can support custom formats as well to support other
	// tools (e.g. cyclondex-json, custom-format-for-tool-ABC, etc.).
	// When creating VMClarity JSON output, use types library from VMClarity API to construct the output.
	OutputFormat ConfigOutputFormat `json:"outputFormat" validate:"required,oneof=vmclarity-json"`

	// TimeoutSeconds The maximum time in seconds that a scan started from this config
	// should run for before being automatically aborted.
	TimeoutSeconds int `json:"timeoutSeconds" validate:"required,gt=0"`
}

// ConfigOutputFormat The format in which the scanner plugin should store it's findings.
// To ensure operability with VMClarity API, the format must be one of enum values.
// However, the scanner can support custom formats as well to support other
// tools (e.g. cyclondex-json, custom-format-for-tool-ABC, etc.).
// When creating VMClarity JSON output, use types library from VMClarity API to construct the output.
type ConfigOutputFormat string

// ErrorResponse An object that is returned for a failed API request.
type ErrorResponse struct {
	Message *string `json:"message,omitempty"`
}

// Metadata Describes the scanner plugin.
type Metadata struct {
	ApiVersion *string `json:"apiVersion,omitempty"`
}

// Status defines model for Status.
type Status struct {
	// LastTransitionTime Last date time when the status has changed.
	LastTransitionTime time.Time `json:"lastTransitionTime"`

	// Message Human readable message.
	Message *string `json:"message,omitempty"`

	// State Describes the status of scanner.
	// | Status         | Description                                                   |
	// | -------------- | ------------------------------------------------------------- |
	// | NotReady       | Initial state when the scanner container starts               |
	// | Ready          | Scanner setup is complete and it is ready to receive requests |
	// | Running        | Scanner config was received and the scanner is running        |
	// | Failed         | Scanner failed                                                |
	// | Done           | Scanner is completed successfully                             |
	State StatusState `json:"state"`
}

// StatusState Describes the status of scanner.
// | Status         | Description                                                   |
// | -------------- | ------------------------------------------------------------- |
// | NotReady       | Initial state when the scanner container starts               |
// | Ready          | Scanner setup is complete and it is ready to receive requests |
// | Running        | Scanner config was received and the scanner is running        |
// | Failed         | Scanner failed                                                |
// | Done           | Scanner is completed successfully                             |
type StatusState string

// PostConfigJSONRequestBody defines body for PostConfig for application/json ContentType.
type PostConfigJSONRequestBody = Config
