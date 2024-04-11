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
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Checkmarx/kics/pkg/model"
	"github.com/Checkmarx/kics/pkg/printer"
	"github.com/Checkmarx/kics/pkg/progress"
	"github.com/Checkmarx/kics/pkg/scan"
	"github.com/hashicorp/hcl/v2/hclsimple"
	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/plugins/sdk/cmd/run"
	"github.com/openclarity/vmclarity/plugins/sdk/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var mapKICSSeverity = map[model.Severity]apitypes.MisconfigurationSeverity{
	model.SeverityHigh:   apitypes.MisconfigurationHighSeverity,
	model.SeverityMedium: apitypes.MisconfigurationMediumSeverity,
	model.SeverityLow:    apitypes.MisconfigurationLowSeverity,
	model.SeverityInfo:   apitypes.MisconfigurationInfoSeverity,
	model.SeverityTrace:  apitypes.MisconfigurationInfoSeverity,
}

//nolint:containedctx
type KICSScanner struct {
	healthz bool
	status  *types.Status
	cancel  context.CancelFunc
}

type ScanParametersConfig struct {
	PreviewLines     int      `json:"preview-lines" yaml:"preview-lines" toml:"preview-lines" hcl:"preview-lines"`
	Platform         []string `json:"platform" yaml:"platform" toml:"platform" hcl:"platform"`
	MaxFileSizeFlag  int      `json:"max-file-size-flag" yaml:"max-file-size-flag" toml:"max-file-size-flag" hcl:"max-file-size-flag"`
	DisableSecrets   bool     `json:"disable-secrets" yaml:"disable-secrets" toml:"disable-secrets" hcl:"disable-secrets"`
	QueryExecTimeout int      `json:"query-exec-timeout" yaml:"query-exec-timeout" toml:"query-exec-timeout" hcl:"query-exec-timeout"`
	Silent           bool     `json:"silent" yaml:"silent" toml:"silent" hcl:"silent"`
	Minimal          bool     `json:"minimal" yaml:"minimal" toml:"minimal" hcl:"minimal"`
}

func (s *KICSScanner) Healthz() bool {
	return s.healthz
}

func (s *KICSScanner) Start(config *types.Config) {
	log.Infof("Starting scanner with config: %+v\n", config)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.TimeoutSeconds)*time.Second)
		s.cancel = cancel
		defer cancel()

		log.Infof("Scanner is running...")
		s.SetStatus(types.NewScannerStatus(types.Running, types.Ptr("Scanner is running...")))
		tmp := os.TempDir()

		clientConfig, err := s.createScanParametersConfig(config.File)
		if err != nil {
			log.Errorf("Failed to parse config file: %v", err)
			s.SetStatus(types.NewScannerStatus(types.Failed, types.Ptr(fmt.Errorf("failed to parse config file: %w", err).Error())))
			return
		}

		c, err := scan.NewClient(
			&scan.Parameters{
				Path:             []string{config.InputDir},
				QueriesPath:      []string{"../../../queries"},
				PreviewLines:     clientConfig.PreviewLines,
				Platform:         clientConfig.Platform,
				OutputPath:       tmp,
				MaxFileSizeFlag:  clientConfig.MaxFileSizeFlag,
				DisableSecrets:   clientConfig.DisableSecrets,
				QueryExecTimeout: clientConfig.QueryExecTimeout,
				OutputName:       "kics",
			},
			&progress.PbBuilder{Silent: clientConfig.Silent},
			printer.NewPrinter(clientConfig.Minimal), //nolint:forbidigo
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

		err = s.formatOutput(tmp, config.OutputDir, config.OutputFormat)
		if err != nil {
			log.Errorf("Failed to format KICS output: %v", err)
			s.SetStatus(types.NewScannerStatus(types.Failed, types.Ptr(fmt.Errorf("failed to format KICS output: %w", err).Error())))
			return
		}

		log.Infof("Scanner finished running.")
		s.SetStatus(types.NewScannerStatus(types.Done, types.Ptr("Scanner finished running.")))
	}()
}

//nolint:gomnd
func (s *KICSScanner) createScanParametersConfig(configPath *string) (*ScanParametersConfig, error) {
	config := ScanParametersConfig{
		PreviewLines:     3,
		Platform:         []string{"Ansible", "CloudFormation", "Common", "Crossplane", "Dockerfile", "DockerCompose", "Knative", "Kubernetes", "OpenAPI", "Terraform", "AzureResourceManager", "GRPC", "GoogleDeploymentManager", "Buildah", "Pulumi", "ServerlessFW", "CICD"},
		MaxFileSizeFlag:  100,
		DisableSecrets:   true,
		QueryExecTimeout: 60,
		Silent:           true,
		Minimal:          true,
	}

	if configPath == nil {
		return &config, nil
	}

	file, err := os.Open(filepath.Clean(*configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	ext := filepath.Ext(*configPath)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	switch ext {
	case ".json":
		if err := json.Unmarshal(bytes, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON config: %w", err)
		} else {
			return &config, nil
		}

	case ".yaml", ".yml":
		if err := yaml.Unmarshal(bytes, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML config: %w", err)
		} else {
			return &config, nil
		}

	case ".toml":
		if _, err := toml.Decode(string(bytes), &config); err != nil {
			return nil, fmt.Errorf("failed to decode TOML config: %w", err)
		} else {
			return &config, nil
		}

	case ".hcl":
		err := hclsimple.DecodeFile(*configPath, nil, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to decode HCL config: %w", err)
		} else {
			return &config, nil
		}

	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}
}

func (s *KICSScanner) formatOutput(tmp, outputDir string, outputFormat types.ConfigOutputFormat) error {
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

	var result []apitypes.Misconfiguration
	for _, q := range summary.Queries {
		for _, file := range q.Files {
			result = append(result, apitypes.Misconfiguration{
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

	var jsonData []byte
	switch outputFormat {
	case types.VMClarityJSON:
		jsonData, err = json.MarshalIndent(apitypes.PluginOutput{Misconfigurations: &result}, "", "    ")
		if err != nil {
			return fmt.Errorf("failed to marshal kics.json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}

	file, err = os.Create(outputDir + "/kics-formatted.json")
	if err != nil {
		return fmt.Errorf("failed to create kics-formatted.json: %w", err)
	}
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

func (s *KICSScanner) Stop(_ int) {
	go func() {
		if s.cancel != nil {
			s.cancel()
		}
	}()
}

func main() {
	k := &KICSScanner{
		healthz: true,
		status:  types.NewScannerStatus(types.Ready, types.Ptr("Starting scanner...")),
	}

	run.Run(k)
}
