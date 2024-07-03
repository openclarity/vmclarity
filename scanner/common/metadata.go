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

package common

import (
	"time"
)

type ScanInputMetadata struct {
	ScannerName string    `json:"scanner_name" yaml:"scanner_name" mapstructure:"scanner_name"`
	InputType   InputType `json:"input_type" yaml:"input_type" mapstructure:"input_type"`
	InputPath   string    `json:"input_path" yaml:"input_path" mapstructure:"input_path"`
	InputSize   int64     `json:"input_size" yaml:"input_size" mapstructure:"input_size"`
	StartTime   time.Time `json:"start_time" yaml:"start_time" mapstructure:"start_time"`
	EndTime     time.Time `json:"end_time" yaml:"end_time" mapstructure:"end_time"`
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

type ScanMetadata struct {
	Inputs    []ScanInputMetadata `json:"inputs" yaml:"inputs" mapstructure:"inputs"`
	StartTime time.Time           `json:"start_time" yaml:"start_time" mapstructure:"start_time"`
	EndTime   time.Time           `json:"end_time" yaml:"end_time" mapstructure:"end_time"`
}

func (s *ScanMetadata) Merge(meta ScanInputMetadata) {
	s.Inputs = append(s.Inputs, meta)

	if s.StartTime.IsZero() || s.StartTime.After(meta.StartTime) {
		s.StartTime = meta.StartTime
	}

	if s.EndTime.IsZero() || s.EndTime.Before(meta.EndTime) {
		s.EndTime = meta.EndTime
	}
}
