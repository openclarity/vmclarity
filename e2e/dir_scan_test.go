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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
)

var _ = ginkgo.Describe("Running a basic scan (only SBOM)", func() {
	reportFailedConfig := ReportFailedConfig{
		services: []string{"orchestrator"},
	}

	ginkgo.Context("which scans a docker container", func() {
		ginkgo.It("should finish successfully", func(ctx ginkgo.SpecContext) {
			var assets *apitypes.Assets
			var err error

			dir := "/tmp/testdir"
			cleanup, err := createTestDir(dir)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer cleanup()

			assetType := apitypes.AssetType{}
			err = assetType.FromDirInfo(apitypes.DirInfo{
				ObjectType: "DirInfo",
				DirName:    to.Ptr("test"),
				Location:   to.Ptr(dir),
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			ginkgo.By("add dir asset")
			_, err = client.PostAsset(
				ctx,
				apitypes.Asset{
					AssetInfo: &assetType,
					FirstSeen: to.Ptr(time.Now()),
				},
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			scope := "assetInfo/objectType eq 'DirInfo' and assetInfo/dirName eq 'test'"

			ginkgo.By("waiting until test asset is found")
			reportFailedConfig.objects = append(
				reportFailedConfig.objects,
				APIObject{"asset", scope},
			)
			assetsParams := apitypes.GetAssetsParams{
				Filter: to.Ptr(scope),
			}
			gomega.Eventually(func() bool {
				assets, err = client.GetAssets(ctx, assetsParams)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				return len(*assets.Items) == 1
			}, DefaultTimeout, time.Second).Should(gomega.BeTrue())

			ginkgo.By("applying a scan configuration")
			apiScanConfig, err := client.PostScanConfig(
				ctx,
				GetCustomScanConfig(
					&apitypes.ScanFamiliesConfig{
						Sbom: &apitypes.SBOMConfig{
							Enabled: to.Ptr(true),
						},
					},
					scope,
					600,
				))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			reportFailedConfig.objects = append(
				reportFailedConfig.objects,
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
					reportFailedConfig.objects = append(
						reportFailedConfig.objects,
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
			}, time.Second*120, time.Second).Should(gomega.BeTrue())

			ginkgo.By("verifying that at least one package was found")
			gomega.Eventually(func() bool {
				totalPackages := (*scans.Items)[0].Summary.TotalPackages
				return *totalPackages > 0
			}, DefaultTimeout, time.Second).Should(gomega.BeTrue())
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
func createTestDir(dir string) (func(), error) {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return nil, fmt.Errorf("failed to create test dir: %w", err)
	}

	// Create a test file to print a UUID
	testFile := filepath.Join(dir, "test.go")
	testScript := "package main\n\nimport (\n\"fmt\"\n\"github.com/google/uuid\"\n)\n\nfunc main() {fmt.Println(uuid.New().String())}\n"
	err = os.WriteFile(testFile, []byte(testScript), 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to create test file: %w", err)
	}

	// Create a mod file to include the UUID package
	testModFile := filepath.Join(dir, "go.mod")
	testModScript := "module test\n\ngo 1.21.7\n\nrequire github.com/google/uuid v1.6.0\n"
	err = os.WriteFile(testModFile, []byte(testModScript), 0o600)
	if err != nil {
		return nil, fmt.Errorf("failed to create test mod file: %w", err)
	}

	return func() {
		os.RemoveAll("/tmp/testdir")
	}, nil
}
