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

package runtimehandler

import (
	"context"
	"io"

	plugintypes "github.com/openclarity/vmclarity/plugins/sdk-go/types"
)

const (
	RemoteScanInputDirOverride   = "/input"           // path in container where the input dir will be mounted
	RemoteScanResultFileOverride = "/tmp/result.json" // path in container where the output result file should be saved by scanner
)

// PluginRuntimeHandler ensures lifecycle and execution consistency
// across plugin runtime implementations.
type PluginRuntimeHandler interface {
	Start(ctx context.Context) error
	Ready() (bool, error)
	GetPluginServerEndpoint(ctx context.Context) (string, error)
	GetOutputFilePath(ctx context.Context) (string, error)
	Logs(ctx context.Context) (io.ReadCloser, error)
	Result(ctx context.Context) (io.ReadCloser, error)
	Remove(ctx context.Context) error
}

// WithOverrides should be used in runner.Start to use valid paths in plugin
// container due to differences in fs mappings.
func WithOverrides(c plugintypes.Config) plugintypes.Config {
	return plugintypes.Config{
		InputDir:       RemoteScanInputDirOverride,
		OutputFile:     c.OutputFile,
		ScannerConfig:  c.ScannerConfig,
		TimeoutSeconds: c.TimeoutSeconds,
	}
}
