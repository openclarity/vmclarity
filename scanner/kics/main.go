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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Checkmarx/kics/pkg/model"
	"github.com/Checkmarx/kics/pkg/printer"
	"github.com/Checkmarx/kics/pkg/progress"
	"github.com/Checkmarx/kics/pkg/scan"
	"github.com/openclarity/vmclarity/scanner/plugin/cmd/run"
	"github.com/openclarity/vmclarity/scanner/types"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"time"
)

var mapKICSSeverity = map[model.Severity]types.MisconfigurationSeverity{
	model.SeverityHigh:   types.MisconfigurationHighSeverity,
	model.SeverityMedium: types.MisconfigurationMediumSeverity,
	model.SeverityLow:    types.MisconfigurationLowSeverity,
	model.SeverityInfo:   types.MisconfigurationInfoSeverity,
	model.SeverityTrace:  types.MisconfigurationInfoSeverity,
}

type KICSScanner struct {
	healthz bool
	status  *types.Status
	ctx     context.Context
}

func (s *KICSScanner) Healthz() bool {
	return s.healthz
}

func (s *KICSScanner) Start(config *types.Config) error {
	log.Infof("Starting scanner with config: %+v\n", config)

	go func() {
		ctx, cancel := context.WithTimeout(s.ctx, time.Duration(config.TimeoutSeconds)*time.Second)
		defer cancel()

		log.Infof("Scanner is running...")
		s.SetStatus(types.NewScannerStatus(types.Running, types.Ptr("Scanner is running...")))
		tmp := os.TempDir()

		c, err := scan.NewClient(
			&scan.Parameters{
				Path:             []string{config.InputDir},
				QueriesPath:      []string{"../../../queries"},
				PreviewLines:     3,
				Platform:         []string{"OpenAPI"},
				OutputPath:       tmp,
				MaxFileSizeFlag:  100,
				DisableSecrets:   true,
				QueryExecTimeout: 60,
				OutputName:       "kics",
			},
			&progress.PbBuilder{Silent: false},
			printer.NewPrinter(true),
		)
		if err != nil {
			log.Errorf("Failed to create KICS client: %v", err)
			s.SetStatus(types.NewScannerStatus(types.Failed, types.Ptr(fmt.Errorf("failed to create KICS client: %w", err).Error())))
			return
		}

		err = c.PerformScan(ctx)
		if err != nil {
			log.Errorf("Failed to perform KICS scan: %v", err)
			s.SetStatus(types.NewScannerStatus(types.Failed, types.Ptr(fmt.Errorf("failed to perform KICS scan: %w", err).Error())))
			return
		}

		if ctx.Err() != nil {
			log.Errorf("The operation timed out: %v", ctx.Err())
			s.SetStatus(types.NewScannerStatus(types.Failed, types.Ptr(fmt.Errorf("failed due to timeout %w", ctx.Err()).Error())))
			return
		}

		err = s.formatOutput(tmp, config.OutputDir)
		if err != nil {
			log.Errorf("Failed to format KICS output: %v", err)
			s.SetStatus(types.NewScannerStatus(types.Failed, types.Ptr(fmt.Errorf("failed to format KICS output: %w", err).Error())))
			return
		}

		log.Infof("Scanner finished running.")
		s.SetStatus(types.NewScannerStatus(types.Done, types.Ptr("Scanner finished running.")))
	}()

	return nil
}

func (s *KICSScanner) formatOutput(tmp, outputDir string) error {
	file, err := os.Open(tmp + "/kics.json")
	if err != nil {
		return fmt.Errorf("failed to open kics.json: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var summary model.Summary
	err = decoder.Decode(&summary)
	if err != nil {
		return fmt.Errorf("failed to decode kics.json: %w", err)
	}

	var result []types.Misconfiguration
	for _, q := range summary.Queries {
		for _, file := range q.Files {
			result = append(result, types.Misconfiguration{
				ScannerName: types.Ptr("KICS"),
				Id:          types.Ptr(file.SimilarityID),
				Location:    types.Ptr(file.FileName + "#" + strconv.Itoa(file.Line)),
				Category:    types.Ptr(q.Category + ":" + string(file.IssueType)),
				Message:     types.Ptr(file.KeyActualValue),
				Description: types.Ptr(q.Description),
				Remediation: types.Ptr(file.KeyExpectedValue),
				Severity:    types.Ptr(mapKICSSeverity[q.Severity]),
			})
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal kics.json: %w", err)
	}

	file, err = os.Create(outputDir + "/kics-formatted.json")
	defer file.Close()

	_, err = io.WriteString(file, string(jsonData))
	if err != nil {
		return fmt.Errorf("failed to write kics-formatted.json: %w", err)
	}

	return nil
}

func (s *KICSScanner) GetStatus() *types.Status {
	return s.status
}

func (s *KICSScanner) SetStatus(newStatus *types.Status) {
	s.status = types.NewScannerStatus(newStatus.State, newStatus.Message)
}

func main() {
	d := &KICSScanner{
		healthz: true,
		status:  types.NewScannerStatus(types.Ready, types.Ptr("Starting scanner...")),
		ctx:     context.Background(),
	}

	run.Run(d)
}
