// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/utils/image_helper"
)

type Source struct {
	Metadata map[string]string `json:"metadata"`
	Type     string            `json:"type"`
	Name     string            `json:"name"` // path in the case of the Type=dir or file, and userInput in the case of Type=image
	Hash     string            `json:"hash"`
}

type Result struct {
	Metadata             families.ScanMetadata              `json:"Metadata"`
	Source               Source                             `json:"Source"`
	VulnerabilitiesByKey map[VulnerabilityKey]Vulnerability `json:"VulnerabilitiesByKey"`
}

func NewResult() *Result {
	return &Result{
		VulnerabilitiesByKey: make(map[VulnerabilityKey]Vulnerability),
	}
}

func (r *Result) SetHash(hash string) {
	r.Source.Hash = hash
}

func (r *Result) SetName(name string) {
	r.Source.Name = name
}

func (r *Result) SetType(srcType string) {
	r.Source.Type = srcType
}

func (r *Result) SetSource(src Source) {
	r.Source = src
}

func (r *Result) GetSourceImageInfo() (*apitypes.ContainerImageInfo, error) {
	sourceImage := image_helper.ImageInfo{}
	if err := sourceImage.FromMetadata(r.Source.Metadata); err != nil {
		return nil, fmt.Errorf("failed to load source image from metadata: %w", err)
	}

	containerImageInfo, err := sourceImage.ToContainerImageInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to convert container image: %w", err)
	}

	return containerImageInfo, nil
}

// ToSlice returns Result in a slice format and not by key.
func (r *Result) ToSlice() []Vulnerability {
	ret := make([]Vulnerability, 0)
	for _, vulnerabilities := range r.VulnerabilitiesByKey {
		ret = append(ret, vulnerabilities)
	}

	return ret
}

func (r *Result) Merge(meta families.ScanInputMetadata, result *ScannerResult) {
	r.Metadata.Merge(meta)

	// Skip further merge if scanner result is empty
	if result == nil {
		return
	}

	for _, vulnerability := range result.Vulnerabilities {
		key := NewVulnerabilityKey(vulnerability)

		// look for other vulnerability key in the current merged vulnerabilities list
		if existingVulnerability, ok := r.VulnerabilitiesByKey[key]; !ok {
			log.Debugf("Adding new vulnerability results from %v. key=%v", result.Scanner, key)
			r.VulnerabilitiesByKey[key] = vulnerability
		} else {
			r.VulnerabilitiesByKey[key] = handleVulnerabilityWithExistingKey(existingVulnerability, vulnerability)
		}
	}

	// TODO: what should we do with other.Source
	// Set Source only once
	if r.Source.Type == "" {
		r.Source = result.Source
	}
}
