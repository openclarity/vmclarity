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

package families

import (
	"context"
	"testing"
	"time"

	misconfigurationTypes "github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	"github.com/openclarity/vmclarity/scanner/families/types"
	"github.com/openclarity/vmclarity/scanner/utils"
	"github.com/stretchr/testify/assert"
)

type familyNotifierSpy struct {
	Results []FamilyResult
}

func (n *familyNotifierSpy) FamilyStarted(context.Context, types.FamilyType) error {
	return nil
}

func (n *familyNotifierSpy) FamilyFinished(_ context.Context, res FamilyResult) error {
	n.Results = append(n.Results, res)

	return nil
}

func TestManagerRunTimeout(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		wantErr    error
		wantResult interface{}
	}{
		{
			name: "Run with misconfiguration family should timeout",
			config: &Config{
				Misconfiguration: misconfigurationTypes.Config{
					Enabled:      true,
					ScannersList: []string{"fake"},
					Inputs: []types.Input{
						{
							Input:     "./",
							InputType: string(utils.ROOTFS),
						},
					},
				},
			},
			wantErr:    context.DeadlineExceeded,
			wantResult: nil,
		},
	}
	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			manager := New(ttp.config)
			notifier := &familyNotifierSpy{}

			ctx, cancel := context.WithTimeout(context.Background(), -time.Nanosecond)
			defer cancel()

			manager.Run(ctx, notifier)
			if ttp.wantErr != nil {
				assert.EqualError(t, ttp.wantErr, ctx.Err().Error())
			}

			for _, res := range notifier.Results {
				if ttp.wantErr != nil {
					assert.Error(t, res.Err, "expected error in family result")
				}
			}
		})
	}
}

// func TestManagerRunPlugin(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		config     *Config
// 		wantErr    error
// 		wantResult interface{}
// 	}{
// 		{
// 			name: "KICS generates raw results",
// 			config: &Config{
// 				Plugins: plugins.Config{
// 					Enabled:      true,
// 					ScannersList: []string{"kics"},
// 					Inputs: []types.Input{
// 						{
// 							Input:     getAbsPathOfTestdata(t, "../../e2e/testdata"),
// 							InputType: string(utils.ROOTFS),
// 						},
// 					},
// 					ScannersConfig: &common.ScannersConfig{
// 						"kics": config.Config{
// 							Name:          "kics",
// 							ImageName:     "ghcr.io/openclarity/vmclarity-plugin-kics:latest",
// 							InputDir:      "",
// 							ScannerConfig: "",
// 						},
// 					},
// 				},
// 			},
// 			wantErr:    nil,
// 			wantResult: float64(23),
// 		},
// 	}
// 	for _, tt := range tests {
// 		ttp := tt
// 		t.Run(ttp.name, func(t *testing.T) {
// 			manager := New(ttp.config)
// 			notifier := &familyNotifierSpy{}

// 			ctx := context.Background()
// 			manager.Run(ctx, notifier)
// 			if ttp.wantErr != nil {
// 				assert.EqualError(t, ttp.wantErr, ctx.Err().Error())
// 			}

// 			for _, res := range notifier.Results {
// 				if ttp.wantErr != nil {
// 					assert.Error(t, res.Err, "expected error in family result")
// 				}

// 				if ttp.wantResult != nil {
// 					fmt.Println(res.Result.(*plugins.Results).RawData["kics"]["total_counter"])
// 					assert.Equal(t, ttp.wantResult, res.Result.(*plugins.Results).RawData["kics"]["total_counter"])
// 				}
// 			}
// 		})
// 	}
// }

// func getAbsPathOfTestdata(t *testing.T, path string) string {
// 	absPath, err := filepath.Abs(path)
// 	if err != nil {
// 		assert.NoError(t, err, "failed to get absolute path")
// 	}

// 	return absPath
// }
