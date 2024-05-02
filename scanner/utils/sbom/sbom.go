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
	"github.com/openclarity/vmclarity/scanner/scanner"
	"github.com/openclarity/vmclarity/scanner/utils/cyclonedx_helper"
)

type CDX struct {
	BOM *cdx.BOM
}

func NewCDX(inputSBOMFile string) (*CDX, error) {
	cdxBOM, err := converter.GetCycloneDXSBOMFromFile(inputSBOMFile)
	if err != nil {
		return nil, converter.ErrFailedToGetCycloneDXSBOM
	}

	return &CDX{
		BOM: cdxBOM,
	}, nil
}

func (c *CDX) GetTargetNameFromSBOM() string {
	return c.BOM.Metadata.Component.Name
}

func (c *CDX) GetHashFromSBOM() (string, error) {
	hash, err := cyclonedx_helper.GetComponentHash(c.BOM.Metadata.Component)
	if err != nil {
		return "", fmt.Errorf("unable to get hash from original SBOM: %w", err)
	}

	return hash, nil
}

func (c *CDX) GetPropertiesFromSBOM() scanner.Metadata {
	return c.GetImageProperties()
}

func (c *CDX) GetImageProperties() scanner.Metadata {
	var imageProperties scanner.Metadata
	for _, property := range *c.BOM.Metadata.Component.Properties {
		switch property.Name {
		case "vmclarity:image:ID":
			imageProperties = append(imageProperties, scanner.Metadata{
				{
					Key:   "ImageID",
					Value: property.Value,
				},
			}...)
		case "vmclarity:image:RepoDigest":
			imageProperties = append(imageProperties, scanner.Metadata{
				{
					Key:   "ImageRepoDigest",
					Value: property.Value,
				},
			}...)
		case "vmclarity:image:Tag":
			imageProperties = append(imageProperties, scanner.Metadata{
				{
					Key:   "ImageTag",
					Value: property.Value,
				},
			}...)
		}
	}

	return imageProperties
}

func SetImageProperties(ID string, RepoDigests []string, Tags []string) []cdx.Property {
	properties := []cdx.Property{}
	if ID != "" {
		properties = append(properties, cdx.Property{
			Name:  "vmclarity:image:ID",
			Value: ID,
		})
	}

	for _, digest := range RepoDigests {
		properties = append(properties, cdx.Property{
			Name:  "vmclarity:image:RepoDigest",
			Value: digest,
		})
	}

	for _, tag := range Tags {
		properties = append(properties, cdx.Property{
			Name:  "vmclarity:image:Tag",
			Value: tag,
		})
	}

	return properties
}
