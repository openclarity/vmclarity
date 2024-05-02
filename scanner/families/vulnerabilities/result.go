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

package vulnerabilities

import (
	"errors"

	"github.com/openclarity/vmclarity/scanner/families/types"
	"github.com/openclarity/vmclarity/scanner/scanner"
	"github.com/openclarity/vmclarity/scanner/utils"
)

type Results struct {
	Metadata      types.Metadata
	MergedResults *scanner.MergedResults
}

func (*Results) IsResults() {}

var imageSourceTypes = map[string]struct{}{
	string(utils.IMAGE):         {},
	string(utils.DOCKERARCHIVE): {},
	string(utils.OCIARCHIVE):    {},
	string(utils.OCIDIR):        {},
}

func (r *Results) GetSourceImageID() (string, error) {
	if r.MergedResults == nil {
		return "", errors.New("missing merged results")
	}

	if _, ok := imageSourceTypes[r.MergedResults.Source.Type]; !ok {
		return "", errors.New("source type is not image")
	}

	for _, prop := range r.MergedResults.Source.Metadata {
		if prop.Key == "ImageID" {
			return prop.Value, nil
		}
	}

	return "", errors.New("missing imageID property")
}

func (r *Results) GetSourceImageRepoDigests() ([]string, error) {
	if r.MergedResults == nil {
		return nil, errors.New("missing merged results")
	}

	if _, ok := imageSourceTypes[r.MergedResults.Source.Type]; !ok {
		return nil, errors.New("source type is not image")
	}

	var repoDigests []string
	for _, prop := range r.MergedResults.Source.Metadata {
		if prop.Key == "ImageRepoDigest" {
			repoDigests = append(repoDigests, prop.Value)
		}
	}

	return repoDigests, nil
}

func (r *Results) GetSourceImageTags() ([]string, error) {
	if r.MergedResults == nil {
		return nil, errors.New("missing merged results")
	}

	if _, ok := imageSourceTypes[r.MergedResults.Source.Type]; !ok {
		return nil, errors.New("source type is not image")
	}

	var tags []string
	for _, prop := range r.MergedResults.Source.Metadata {
		if prop.Key == "ImageTag" {
			tags = append(tags, prop.Value)
		}
	}

	return tags, nil
}
