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

package sbom

import (
	"fmt"

	cdx "github.com/CycloneDX/cyclonedx-go"

	"github.com/openclarity/vmclarity/scanner/converter"
	"github.com/openclarity/vmclarity/scanner/utils/cyclonedx_helper"
)

type CycloneDXSBOM struct {
	BOM *cdx.BOM
}

func NewCycloneDXSBOM(inputSBOMFile string) (*CycloneDXSBOM, error) {
	cdxBOM, err := converter.GetCycloneDXSBOMFromFile(inputSBOMFile)
	if err != nil {
		return nil, converter.ErrFailedToGetCycloneDXSBOM
	}

	return &CycloneDXSBOM{
		BOM: cdxBOM,
	}, nil
}

func (c *CycloneDXSBOM) GetTargetNameFromSBOM() string {
	return c.BOM.Metadata.Component.Name
}

func (c *CycloneDXSBOM) GetHashFromSBOM() (string, error) {
	hash, err := cyclonedx_helper.GetComponentHash(c.BOM.Metadata.Component)
	if err != nil {
		return "", fmt.Errorf("unable to get hash from original SBOM: %w", err)
	}

	return hash, nil
}

func (c *CycloneDXSBOM) GetPropertiesFromSBOM() []cdx.Property {
	return *c.BOM.Metadata.Component.Properties
}
