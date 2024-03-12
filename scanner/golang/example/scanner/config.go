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

package scanner

import (
	"time"

	dockle_config "github.com/Portshift/dockle/config"
	"github.com/openclarity/vmclarity/scanner/types"
)

func createDockleConfig(sourceType types.ScanObjectInputType, input string) *dockle_config.Config {
	dockleConfig := &dockle_config.Config{
		Debug:      true,
		Timeout:    2 * time.Minute,
		LocalImage: true,
	}

	// nolint:exhaustive
	switch sourceType {
	case types.InputTypeDockerArchive:
		dockleConfig.FilePath = input
	default:
		dockleConfig.ImageName = input
	}

	return dockleConfig
}
