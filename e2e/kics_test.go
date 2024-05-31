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

package e2e

import (
	"context"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/plugins"
	"github.com/openclarity/vmclarity/scanner/families/plugins/common"
	"github.com/openclarity/vmclarity/scanner/families/plugins/runner/config"
	"github.com/openclarity/vmclarity/scanner/families/types"
	"github.com/openclarity/vmclarity/scanner/utils"
)

type Notifier struct {
	Results []families.FamilyResult
}

func (n *Notifier) FamilyStarted(context.Context, types.FamilyType) error { return nil }

func (n *Notifier) FamilyFinished(_ context.Context, res families.FamilyResult) error {
	n.Results = append(n.Results, res)

	return nil
}

var _ = ginkgo.Describe("Running KICS scan", func() {
	ginkgo.Context("which scans an openapi.yaml file", func() {
		ginkgo.It("should finish successfully", func(ctx ginkgo.SpecContext) {
			image := os.Getenv("VMCLARITY_E2E_PLUGIN_KICS_IMAGE")
			input, err := filepath.Abs("./testdata")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			notifier := &Notifier{}

			families.New(&families.Config{
				Plugins: plugins.Config{
					Enabled:      true,
					ScannersList: []string{"kics"},
					Inputs: []types.Input{
						{
							Input:     input,
							InputType: string(utils.ROOTFS),
						},
					},
					ScannersConfig: &common.ScannersConfig{
						"kics": config.Config{
							Name:          "kics",
							ImageName:     image,
							InputDir:      "",
							ScannerConfig: "",
						},
					},
				},
			}).Run(context.TODO(), notifier)

			gomega.Eventually(func() bool {
				for _, res := range notifier.Results {
					return gomega.Expect(res.Result.(*plugins.Results).RawData["kics"].(map[string]interface{})["total_counter"]).To(gomega.Equal(float64(23))) &&
						gomega.Expect(len(res.Result.(*plugins.Results).Output)).To(gomega.Equal(23))
				}
				return false
			}, DefaultTimeout, DefaultPeriod).Should(gomega.BeTrue())
		})
	})
})
