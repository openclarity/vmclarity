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

package fake

import (
	"context"
	familiestypes "github.com/openclarity/vmclarity/scanner/families"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
)

const ScannerName = "fake"

type Scanner struct {
	logger *log.Entry
}

func New(_ string, _ types.ScannersConfig, logger *log.Entry) (familiestypes.Scanner[*types.ScannerResult], error) {
	return &Scanner{
		logger: logger.Dup().WithField("scanner", ScannerName),
	}, nil
}

func (a *Scanner) Scan(_ context.Context, _ scannertypes.InputType, _ string) (*types.ScannerResult, error) {
	return &types.ScannerResult{
		ScannerName:       ScannerName,
		Misconfigurations: createFakeMisconfigurationReport(),
	}, nil
}

func createFakeMisconfigurationReport() []types.Misconfiguration {
	return []types.Misconfiguration{
		{
			Location: "/fake",

			Category:    "FAKE",
			ID:          "Test1",
			Description: "Fake test number 1",

			Message:     "Fake test number 1 failed",
			Severity:    types.HighSeverity,
			Remediation: "fix the thing number 1",
		},
		{
			Location: "/fake",

			Category:    "FAKE",
			ID:          "Test2",
			Description: "Fake test number 2",

			Message:     "Fake test number 2 failed",
			Severity:    types.LowSeverity,
			Remediation: "fix the thing number 2",
		},
		{
			Location: "/fake",

			Category:    "FAKE",
			ID:          "Test3",
			Description: "Fake test number 3",

			Message:     "Fake test number 3 failed",
			Severity:    types.MediumSeverity,
			Remediation: "fix the thing number 3",
		},
		{
			Location: "/fake",

			Category:    "FAKE",
			ID:          "Test4",
			Description: "Fake test number 4",

			Message:     "Fake test number 4 failed",
			Severity:    types.HighSeverity,
			Remediation: "fix the thing number 4",
		},
	}
}

func init() {
	types.FactoryRegister(ScannerName, New)
}
