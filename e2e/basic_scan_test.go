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
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
	"github.com/openclarity/vmclarity/testenv/types"
)

var _ = ginkgo.Describe("Running a basic scan (only SBOM)", func() {
	reportFailedConfig := ReportFailedConfig{}
	var imageID string
	var assets *apitypes.Assets

	ginkgo.Context("which scans an asset", func() {
		ginkgo.It("should finish successfully", func(ctx ginkgo.SpecContext) {
			var err error
			ginkgo.By("waiting until test asset is found")
			reportFailedConfig.objects = append(
				reportFailedConfig.objects,
				APIObject{"asset", cfg.TestSuiteParams.Scope},
			)
			assetsParams := apitypes.GetAssetsParams{
				Filter: to.Ptr(cfg.TestSuiteParams.Scope),
			}
			gomega.Eventually(func() bool {
				assets, err = client.GetAssets(ctx, assetsParams)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				return len(*assets.Items) == 1
			}, DefaultTimeout, time.Second).Should(gomega.BeTrue())

			RunSuccessfulScan(ctx, &reportFailedConfig, cfg.TestSuiteParams.Scope)
		})
	})

	reportFailedConfig = ReportFailedConfig{}

	ginkgo.Context("which scans a docker image", func() {
		ginkgo.It("should finish successfully", func(ctx ginkgo.SpecContext) {
			if cfg.TestEnvConfig.Platform != types.EnvironmentTypeDocker {
				ginkgo.By("skipping test because it's not running on docker")
				return
			}

			containerInfo, err := (*assets.Items)[0].AssetInfo.AsContainerInfo()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			imageID = containerInfo.Image.ImageID

			ginkgo.By("waiting until test asset is found")
			filter := fmt.Sprintf("assetInfo/objectType eq 'ContainerImageInfo' and assetInfo/imageID eq '%s'", imageID)
			reportFailedConfig.objects = append(
				reportFailedConfig.objects,
				APIObject{"asset", filter},
			)
			gomega.Eventually(func() bool {
				assets, err := client.GetAssets(ctx, apitypes.GetAssetsParams{
					Filter: to.Ptr(filter),
				})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				return len(*assets.Items) == 1
			}, DefaultTimeout, time.Second).Should(gomega.BeTrue())

			RunSuccessfulScan(ctx, &reportFailedConfig, filter)
		})
	})

	ginkgo.AfterEach(func(ctx ginkgo.SpecContext) {
		if ginkgo.CurrentSpecReport().Failed() {
			reportFailedConfig.startTime = ginkgo.CurrentSpecReport().StartTime
			ReportFailed(ctx, testEnv, client, &reportFailedConfig)
		}
	})
})

// nolint:gomnd
func RunSuccessfulScan(ctx ginkgo.SpecContext, report *ReportFailedConfig, filter string) {
	ginkgo.By("applying a scan configuration")
	apiScanConfig, err := client.PostScanConfig(
		ctx,
		GetCustomScanConfig(
			&apitypes.ScanFamiliesConfig{
				Sbom: &apitypes.SBOMConfig{
					Enabled: to.Ptr(true),
				},
			},
			filter,
			int(cfg.TestSuiteParams.ScanTimeout.Seconds()),
		))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	report.objects = append(
		report.objects,
		APIObject{"scanConfig", fmt.Sprintf("id eq '%s'", *apiScanConfig.Id)},
	)

	ginkgo.By("updating scan configuration to run now")
	updateScanConfig := UpdateScanConfigToStartNow(apiScanConfig)
	err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By("waiting until scan starts")
	scanParams := apitypes.GetScansParams{
		Filter: to.Ptr(fmt.Sprintf(
			"scanConfig/id eq '%s' and status/state ne '%s' and status/state ne '%s'",
			*apiScanConfig.Id,
			apitypes.ScanStatusStateDone,
			apitypes.ScanStatusStateFailed,
		)),
	}
	var scans *apitypes.Scans
	gomega.Eventually(func() bool {
		scans, err = client.GetScans(ctx, scanParams)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		if len(*scans.Items) == 1 {
			report.objects = append(
				report.objects,
				APIObject{"scan", fmt.Sprintf("id eq '%s'", *(*scans.Items)[0].Id)},
			)
			return true
		}
		return false
	}, DefaultTimeout, time.Second).Should(gomega.BeTrue())

	ginkgo.By("waiting until scan state changes to done")
	scanParams = apitypes.GetScansParams{
		Filter: to.Ptr(fmt.Sprintf(
			"scanConfig/id eq '%s' and status/state eq '%s' and status/reason eq '%s'",
			*apiScanConfig.Id,
			apitypes.AssetScanStatusStateDone,
			apitypes.AssetScanStatusReasonSuccess,
		)),
	}
	gomega.Eventually(func() bool {
		scans, err = client.GetScans(ctx, scanParams)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		return len(*scans.Items) == 1
	}, cfg.TestSuiteParams.ScanTimeout, time.Second).Should(gomega.BeTrue())
}
