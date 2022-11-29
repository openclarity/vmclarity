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

package cloudinit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

func GenerateCloudInit(scannerConfig *types.ScannerConfig) (*string, error) {
	vars := make(map[string]interface{})
	// parse the template
	tmpl, _ := template.New("cloud-init").Parse(cloudInitTmpl)

	scannerConfigB, err := json.Marshal(scannerConfig)
	if err != nil {
		return nil, fmt.Errorf("falied to marshal config: %v", err)
	}

	if err != nil {
		return nil, err
	}
	vars["Config"] = bytes.NewBuffer(scannerConfigB).String()
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, vars); err != nil {
		return nil, err
	}

	cloudInit := tpl.String()
	return &cloudInit, nil
}
