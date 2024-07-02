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

package syft

import (
	"context"
	"errors"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/families/sbom/types"

	"github.com/openclarity/vmclarity/scanner/utils"

	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/cataloging"
	"github.com/anchore/syft/syft/format/common/cyclonedxhelpers"
	syftsbom "github.com/anchore/syft/syft/sbom"
	syftsrc "github.com/anchore/syft/syft/source"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/sbom/syft/config"
	"github.com/openclarity/vmclarity/scanner/job_manager"
	"github.com/openclarity/vmclarity/scanner/utils/image_helper"
)

const AnalyzerName = "syft"

type Analyzer struct {
	name       string
	logger     *log.Entry
	config     config.Config
	resultChan chan job_manager.Result
}

func New(_ string, c job_manager.IsConfig, logger *log.Entry, resultChan chan job_manager.Result) job_manager.Job {
	conf := c.(*types.AnalyzersConfig) // nolint:forcetypeassert
	return &Analyzer{
		name:       AnalyzerName,
		logger:     logger.Dup().WithField("analyzer", AnalyzerName),
		config:     conf.Syft,
		resultChan: resultChan,
	}
}

func (a *Analyzer) Run(ctx context.Context, sourceType utils.SourceType, userInput string) error {
	src := utils.CreateSource(sourceType, a.config.LocalImageScan)

	a.logger.Infof("Called %s analyzer on source %s", a.name, src)
	// TODO platform can be defined
	// https://github.com/anchore/syft/blob/b20310eaf847c259beb4fe5128c842bd8aa4d4fc/cmd/syft/cli/options/packages.go#L48
	source, err := syft.GetSource(
		ctx,
		userInput,
		syft.DefaultGetSourceConfig().
			WithSources(src).
			WithRegistryOptions(a.config.GetRegistryOptions()).
			WithExcludeConfig(a.config.GetExcludePaths()),
	)
	if err != nil {
		return fmt.Errorf("failed to create source analyzer=%s: %w", a.name, err)
	}

	go func(ctx context.Context) {
		res := &types.ScannerResult{}

		sbomConfig := syft.DefaultCreateSBOMConfig().
			WithSearchConfig(cataloging.DefaultSearchConfig().WithScope(a.config.GetScope()))

		sbom, err := syft.CreateSBOM(ctx, source, sbomConfig)
		if err != nil {
			a.setError(res, fmt.Errorf("failed to write results: %w", err))
			return
		}

		cdxBom := cyclonedxhelpers.ToFormatModel(*sbom)
		res = types.CreateScannerResult(cdxBom, a.name, userInput, sourceType)

		// Syft uses ManifestDigest to fill version information in the case of an image.
		// We need RepoDigest/ImageID as well which is not set by Syft if we're using cycloneDX output.
		// Get the RepoDigest/ImageID from image metadata and use it as SourceHash in the Result
		// that will be added to the component hash of metadata during the merge.
		switch sourceType {
		case utils.IMAGE, utils.DOCKERARCHIVE, utils.OCIDIR, utils.OCIARCHIVE:
			hash, imageInfo, err := getImageInfo(sbom, userInput)
			if err != nil {
				a.setError(res, fmt.Errorf("failed to get image hash from sbom: %w", err))
				return
			}

			// sync image details to result
			res.AppInfo.SourceHash = hash
			res.AppInfo.SourceMetadata = imageInfo.ToMetadata()

		case utils.SBOM, utils.DIR, utils.ROOTFS, utils.FILE:
			// ignore
		default:
			// ignore
		}

		a.logger.Infof("Sending successful results")
		a.resultChan <- res
	}(ctx)

	return nil
}

func (a *Analyzer) setError(res *types.ScannerResult, err error) {
	res.Error = err
	a.logger.Error(res.Error)
	a.resultChan <- res
}

func getImageInfo(s *syftsbom.SBOM, src string) (string, *image_helper.ImageInfo, error) {
	switch metadata := s.Source.Metadata.(type) {
	case syftsrc.ImageMetadata:
		imageInfo := &image_helper.ImageInfo{
			Name:    src,
			ID:      metadata.ID,
			Tags:    metadata.Tags,
			Digests: metadata.RepoDigests,
		}

		hash, err := imageInfo.GetHashFromRepoDigestsOrImageID()
		if err != nil {
			return "", nil, fmt.Errorf("failed to get image hash from repo digests or image id: %w", err)
		}

		return hash, imageInfo, nil
	default:
		return "", nil, errors.New("failed to get image hash from source metadata")
	}
}
