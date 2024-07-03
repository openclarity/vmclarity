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

package grype

import (
	"context"
	"encoding/json"
	"fmt"
	familiestypes "github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/vulnerabilities/types"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"
	"os"
	"time"

	grype_models "github.com/anchore/grype/grype/presenter/models"
	transport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	grype_client "github.com/openclarity/grype-server/api/client/client"
	grype_client_operations "github.com/openclarity/grype-server/api/client/client/operations"
	grype_client_models "github.com/openclarity/grype-server/api/client/models"
	log "github.com/sirupsen/logrus"

	sbom "github.com/openclarity/vmclarity/scanner/utils/sbom"
)

type RemoteScanner struct {
	logger  *log.Entry
	client  *grype_client.GrypeServer
	timeout time.Duration
}

func newRemoteScanner(config types.ScannersConfig, logger *log.Entry) familiestypes.Scanner[*types.ScannerResult] {
	cfg := grype_client.DefaultTransportConfig().
		WithSchemes(config.Grype.Remote.GrypeServerSchemes).
		WithHost(config.Grype.Remote.GrypeServerAddress)

	return &RemoteScanner{
		logger:  logger.Dup().WithField("scanner", ScannerName).WithField("scanner-mode", "remote"),
		client:  grype_client.New(transport.New(cfg.Host, cfg.BasePath, cfg.Schemes), strfmt.Default),
		timeout: config.Grype.Remote.GrypeServerTimeout,
	}
}

func (s *RemoteScanner) Scan(ctx context.Context, sourceType scannertypes.InputType, userInput string) (*types.ScannerResult, error) {
	// remote-grype supports only SBOM as a source input since it sends the SBOM to a centralized grype server for scanning.
	if sourceType != scannertypes.SBOM {
		s.logger.Infof("Ignoring non SBOM input. type=%v", sourceType)
		return nil, nil
	}

	sbomBytes, err := os.ReadFile(userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	doc, err := s.scanSbomWithGrypeServer(sbomBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to scan sbom with grype server: %w", err)
	}

	bom, err := sbom.NewCycloneDX(userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create CycloneDX SBOM: %w", err)
	}

	targetName := bom.GetTargetNameFromSBOM()
	metadata := bom.GetMetadataFromSBOM()
	hash, err := bom.GetHashFromSBOM()
	if err != nil {
		return nil, fmt.Errorf("failed to get original hash from SBOM: %w", err)
	}

	s.logger.Infof("Sending successful results")
	result := createResults(*doc, targetName, ScannerName, hash, metadata)

	return result, nil
}

func (s *RemoteScanner) scanSbomWithGrypeServer(sbom []byte) (*grype_models.Document, error) {
	params := grype_client_operations.NewPostScanSBOMParams().
		WithBody(&grype_client_models.SBOM{
			Sbom: sbom,
		}).WithTimeout(s.timeout)
	ok, err := s.client.Operations.PostScanSBOM(params)
	if err != nil {
		return nil, fmt.Errorf("failed to send sbom for scan: %w", err)
	}
	doc := grype_models.Document{}

	err = json.Unmarshal(ok.Payload.Vulnerabilities, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vulnerabilities document: %w", err)
	}

	return &doc, nil
}
