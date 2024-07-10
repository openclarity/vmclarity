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

package e2e

import (
	"context"
	"path/filepath"
	"reflect"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/openclarity/vmclarity/scanner"
	scannercommon "github.com/openclarity/vmclarity/scanner/common"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/plugins/runner/config"
	plugintypes "github.com/openclarity/vmclarity/scanner/families/plugins/types"
)

const scannerPluginNameKICS = "kics"

type Notifier struct {
	Results []families.FamilyResult
}

func (n *Notifier) FamilyStarted(context.Context, families.FamilyType) error { return nil }

func (n *Notifier) FamilyFinished(_ context.Context, res families.FamilyResult) error {
	n.Results = append(n.Results, res)

	return nil
}

var _ = ginkgo.Describe("Running KICS scan", func() {
	ginkgo.Context("which scans an openapi.yaml file", func() {
		ginkgo.It("should finish successfully", func(ctx ginkgo.SpecContext) {
			if cfg.TestEnvConfig.Images.PluginKics == "" {
				ginkgo.Skip("KICS plugin image not provided")
			}

			input, err := filepath.Abs("./testdata")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			notifier := &Notifier{}

			errs := scanner.New(&scanner.Config{
				Plugins: plugintypes.Config{
					Enabled:      true,
					ScannersList: []string{scannerPluginNameKICS},
					Inputs: []scannercommon.ScanInput{
						{
							Input:     input,
							InputType: scannercommon.ROOTFS,
						},
					},
					ScannersConfig: plugintypes.ScannersConfig{
						scannerPluginNameKICS: config.Config{
							Name:          scannerPluginNameKICS,
							ImageName:     cfg.TestEnvConfig.Images.PluginKics,
							InputDir:      "",
							ScannerConfig: "",
						},
					},
				},
			}).Run(ctx, notifier)
			gomega.Expect(errs).To(gomega.BeEmpty())

			gomega.Eventually(func() bool {
				if len(notifier.Results) != 1 {
					return false
				}

				results := notifier.Results[0].Result.(*plugintypes.Result)                              // nolint:forcetypeassert
				rawData := results.PluginOutputs[scannerPluginNameKICS].RawJSON.(map[string]interface{}) // nolint:forcetypeassert

				if rawData["total_counter"] != float64(23) {
					return false
				}

				if len(results.Findings) != 23 {
					return false
				}

				return true
			}, DefaultTimeout, DefaultPeriod).Should(gomega.BeTrue())
		})
	})
})

var _ = ginkgo.Describe("Running a KICS scan", func() {
	ginkgo.Context("which scans an openapi.yaml file and has report-formats set to sarif", func() {
		ginkgo.It("should finish successfully, and output both JSON and Sarif format as well as VMClarity output", func(ctx ginkgo.SpecContext) {
			if cfg.TestEnvConfig.Images.PluginKics == "" {
				ginkgo.Skip("KICS plugin image not provided")
			}

			input, err := filepath.Abs("./testdata")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			notifier := &Notifier{}

			errs := scanner.New(&scanner.Config{
				Plugins: plugintypes.Config{
					Enabled:      true,
					ScannersList: []string{scannerPluginNameKICS},
					Inputs: []scannercommon.ScanInput{
						{
							Input:     input,
							InputType: scannercommon.ROOTFS,
						},
					},
					ScannersConfig: plugintypes.ScannersConfig{
						scannerPluginNameKICS: config.Config{
							Name:          scannerPluginNameKICS,
							ImageName:     cfg.TestEnvConfig.Images.PluginKics,
							InputDir:      "",
							ScannerConfig: "{\"report-formats\": [\"sarif\"]}",
						},
					},
				},
			}).Run(ctx, notifier)
			gomega.Expect(errs).To(gomega.BeEmpty())

			gomega.Eventually(func() bool {
				results := notifier.Results[0].Result.(*plugintypes.Result).PluginOutputs[scannerPluginNameKICS] // nolint:forcetypeassert

				isEmptyFuncs := []func() bool{
					func() bool { return isEmpty(results.RawJSON) },
					func() bool { return isEmpty(results.RawSarif) },
					func() bool { return isEmpty(results.Vmclarity) },
				}

				for _, f := range isEmptyFuncs {
					if f() {
						return false
					}
				}

				return true
			}, DefaultTimeout, DefaultPeriod).Should(gomega.BeTrue())
		})
	})
})

func isEmpty(x interface{}) bool {
	if x == nil {
		return true
	}

	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
