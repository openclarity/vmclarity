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

package end_to_end_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/end_to_end_test/helpers"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
	"time"
)

var _ = Describe("Detecting scan failures", func() {

	Context("when a scan stops without assets to scan", func() {
		It("should detect failure reason successfully", func(ctx SpecContext) {
			By("applying a scan configuration with not existing label")
			apiScanConfig, err := client.PostScanConfig(
				ctx,
				helpers.GetCustomScanConfig(
					&helpers.DefaultScanFamiliesConfig,
					"contains(assetInfo.tags, '{\"key\":\"notexisting\",\"value\":\"label\"}')",
					600,
				))
			Expect(err).NotTo(HaveOccurred())

			By("updating scan configuration to run now")
			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			By("waiting until scan state changes to failed with nothing to scan as state reason")
			params := models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state eq '%s' and stateReason eq '%s'",
					*apiScanConfig.Id,
					models.ScanStateDone,
					models.ScanRelationshipStateReasonNothingToScan,
				)),
			}
			var scans *models.Scans
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, helpers.DefaultTimeout, time.Second).Should(BeTrue())
		})
	})

	Context("when a scan stops with timeout", func() {
		It("should detect failure reason successfully", func(ctx SpecContext) {
			By("applying a scan configuration with short timeout")
			apiScanConfig, err := client.PostScanConfig(
				ctx,
				helpers.GetCustomScanConfig(
					&helpers.DefaultScanFamiliesConfig,
					helpers.DefaultScope,
					2,
				))
			Expect(err).NotTo(HaveOccurred())

			By("updating scan configuration to run now")
			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			By("waiting until scan state changes to failed with timed out as state reason")
			params := models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state eq '%s' and stateReason eq '%s'",
					*apiScanConfig.Id,
					models.ScanStateFailed,
					models.ScanStateReasonTimedOut,
				)),
			}
			var scans *models.Scans
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, helpers.DefaultTimeout, time.Second).Should(BeTrue())
		})
	})
})
