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

package common

import "time"

type ScanInput struct {
	// StripPathFromResult overrides global StripInputPaths value
	StripPathFromResult *bool     `yaml:"strip_path_from_result" mapstructure:"strip_path_from_result"`
	Input               string    `yaml:"input" mapstructure:"input"`
	InputType           InputType `yaml:"input_type" mapstructure:"input_type"`
}

type ScanInputMetadata struct {
	ScannerName string
	InputType   InputType
	InputPath   string
	InputSize   int64
	StartTime   time.Time
	EndTime     time.Time
}

func NewScanInputMetadata(scannerName string, startTime, endTime time.Time, inputSize int64, input ScanInput) ScanInputMetadata {
	return ScanInputMetadata{
		ScannerName: scannerName,
		InputType:   input.InputType,
		InputPath:   input.Input,
		InputSize:   inputSize,
		StartTime:   startTime,
		EndTime:     endTime,
	}
}
