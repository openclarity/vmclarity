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
	"github.com/labstack/echo/v4"
	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/scanner/plugin"
	"github.com/openclarity/vmclarity/scanner/plugin/cmd/run"
	"github.com/openclarity/vmclarity/scanner/types"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

var mapKICSSeverity = map[model.Severity]apitypes.MisconfigurationSeverity{
	model.SeverityHigh:   apitypes.MisconfigurationHighSeverity,
	model.SeverityMedium: apitypes.MisconfigurationMediumSeverity,
	model.SeverityLow:    apitypes.MisconfigurationLowSeverity,
	model.SeverityInfo:   apitypes.MisconfigurationInfoSeverity,
	model.SeverityTrace:  apitypes.MisconfigurationInfoSeverity,
}

type KICSScanner struct {
	healthz bool
	status  *types.Status
}

func (d *KICSScanner) Healthz() bool {
	return d.healthz
}

func (d *KICSScanner) Start(ctx echo.Context, config *types.Config) error {
	log.Infof("Starting scanner with config: %+v\n", config)

	go func() {
		d.SetStatus(types.NewScannerStatus(types.Running, plugin.PointerTo("Scanner is running...")))

		c, err := scan.NewClient(
			&scan.Parameters{
				Path:             []string{config.InputDir},
				QueriesPath:      []string{"../../../queries"},
				PreviewLines:     3,
				Platform:         []string{"OpenAPI"},
				OutputPath:       config.OutputDir,
				MaxFileSizeFlag:  100,
				DisableSecrets:   true,
				QueryExecTimeout: 60,
				OutputName:       "kics",
			},
			&progress.PbBuilder{Silent: false},
			printer.NewPrinter(true),
		)
		if err != nil {
			d.SetStatus(types.NewScannerStatus(types.Failed, plugin.PointerTo(fmt.Sprintf("Failed to initialize scanner: %v", err))))
			fmt.Println(err)
			return
		}

		err = c.PerformScan(context.Background())
		if err != nil {
			d.SetStatus(types.NewScannerStatus(types.Failed, plugin.PointerTo(fmt.Sprintf("Failed to perform scan: %v", err))))
			fmt.Println(err)
			return
		}

		file, err := os.Open(config.OutputDir + "/kics.json")
		if err != nil {
			d.SetStatus(types.NewScannerStatus(types.Failed, plugin.PointerTo(fmt.Sprintf("Failed to open file: %v", err))))
			fmt.Println(err)
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		var summary model.Summary
		err = decoder.Decode(&summary)
		if err != nil {
			d.SetStatus(types.NewScannerStatus(types.Failed, plugin.PointerTo(fmt.Sprintf("Failed to decode json: %v", err))))
			fmt.Println(err)
			return
		}

		var result []apitypes.Misconfiguration
		for _, q := range summary.Queries {
			for _, file := range q.Files {
				result = append(result, apitypes.Misconfiguration{
					ScannerName: plugin.PointerTo("KICS"),
					Id:          plugin.PointerTo(file.SimilarityID),
					Location:    plugin.PointerTo(file.FileName + "#" + strconv.Itoa(file.Line)),
					Category:    plugin.PointerTo(q.Category + ":" + string(file.IssueType)),
					Message:     plugin.PointerTo(file.KeyActualValue),
					Description: plugin.PointerTo(q.Description),
					Remediation: plugin.PointerTo(file.KeyExpectedValue),
					Severity:    plugin.PointerTo(mapKICSSeverity[q.Severity]),
				})
			}
		}

		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			d.SetStatus(types.NewScannerStatus(types.Failed, plugin.PointerTo(fmt.Sprintf("Failed to marshal json: %v", err))))
			fmt.Println(err)
			return
		}

		file, err = os.Create(config.OutputDir + "/kics-formatted.json")
		defer file.Close()

		_, err = io.WriteString(file, string(jsonData))
		if err != nil {
			d.SetStatus(types.NewScannerStatus(types.Failed, plugin.PointerTo(fmt.Sprintf("Failed to write json: %v", err))))
			fmt.Println(err)
			return
		}

		d.SetStatus(types.NewScannerStatus(types.Done, plugin.PointerTo("Scanner finished running.")))
	}()

	return nil
}

func (d *KICSScanner) GetStatus() *types.Status {
	return d.status
}

func (d *KICSScanner) SetStatus(s *types.Status) {
	d.status = types.NewScannerStatus(s.State, s.Message)
}

func main() {
	d := &KICSScanner{
		healthz: true,
		status:  types.NewScannerStatus(types.Ready, plugin.PointerTo("Starting scanner...")),
	}

	run.Run(d)
}
