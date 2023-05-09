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

package state

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/shared/pkg/families/types"
)

type LocalState struct{}

func (l *LocalState) WaitForVolumeAttachment(context.Context) error {
	return nil
}

func (l *LocalState) MarkInProgress(context.Context) error {
	log.Info("Scanning is in progress")
	return nil
}

func (l *LocalState) MarkFamilyScanInProgress(ctx context.Context, familyType types.FamilyType) error {
	var err error
	switch familyType {
	case types.SBOM:
		err = l.markSBOMScanInProgress(ctx)
	case types.Vulnerabilities:
		err = l.markVulnerabilitiesScanInProgress(ctx)
	case types.Secrets:
		err = l.markSecretsScanInProgress(ctx)
	case types.Exploits:
		err = l.markExploitsScanInProgress(ctx)
	case types.Misconfiguration:
		err = l.markMisconfigurationsScanInProgress(ctx)
	case types.Rootkits:
		err = l.markRootkitsScanInProgress(ctx)
	case types.Malware:
		err = l.markRootkitsScanInProgress(ctx)
	}
	return err
}

func (l *LocalState) markExploitsScanInProgress(context.Context) error {
	log.Info("Exploits scan is in progress")
	return nil
}

func (l *LocalState) markSecretsScanInProgress(context.Context) error {
	log.Info("Secrets scan is in progress")
	return nil
}

func (l *LocalState) markSBOMScanInProgress(context.Context) error {
	log.Info("SBOM scan is in progress")
	return nil
}

func (l *LocalState) markVulnerabilitiesScanInProgress(context.Context) error {
	log.Info("Vulnerabilities scan is in progress")
	return nil
}

func (l *LocalState) markMalwareScanInProgress(context.Context) error {
	log.Info("Malware scan is in progress")
	return nil
}

func (l *LocalState) markMisconfigurationsScanInProgress(context.Context) error {
	log.Info("Misconfiguration scan is in progress")
	return nil
}

func (l *LocalState) markRootkitsScanInProgress(context.Context) error {
	log.Info("Rootkit scan is in progress")
	return nil
}

func (l *LocalState) MarkDone(_ context.Context, errs []error) error {
	if len(errs) > 0 {
		log.Errorf("scan has been completed with errors: %v", errs)
		return nil
	}
	log.Info("Scan has been completed")
	return nil
}

func (l *LocalState) IsAborted(context.Context) (bool, error) {
	return false, nil
}

func NewLocalState() (*LocalState, error) {
	return &LocalState{}, nil
}
