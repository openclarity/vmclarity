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

package initiator

import (
	"context"
	"fmt"

	cliconfig "github.com/openclarity/vmclarity/cli/pkg/config"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"github.com/openclarity/vmclarity/shared/pkg/families"
)

type Config struct {
	client         *backendclient.BackendClient
	fmConfig       *families.Config
	scanConfigID   string
	scanConfigName string
	input          string
	asset          cliconfig.Asset
}

func CreateConfig(
	client *backendclient.BackendClient,
	fmConfig *families.Config,
	scanConfigID, scanConfigName, input string,
	asset cliconfig.Asset,
) Config {
	return Config{
		client:         client,
		fmConfig:       fmConfig,
		scanConfigID:   scanConfigID,
		scanConfigName: scanConfigName,
		input:          input,
		asset:          asset,
	}
}

// InitResults creates VMClarityInitiator and init Results.
// The function is returns the scanID and scanResultID that required for the export.
func InitResults(ctx context.Context, standaloneInitiatorConfig Config) (string, string, error) {
	i, err := newVMClarityInitiator(standaloneInitiatorConfig)
	if err != nil {
		return "", "", fmt.Errorf("failed to create VMClarity initiator: %w", err)
	}
	scanID, scanResultID, err := i.initResults(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to init scan result: %w", err)
	}

	return scanID, scanResultID, nil
}
