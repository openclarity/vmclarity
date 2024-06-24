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
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
)

var _ = ginkgo.Describe("Detecting finding summary updates", func() {
	reportFailedConfig := ReportFailedConfig{}

	ginkgo.Context("which updates total high vulnerabilities count", func() {
		ginkgo.It("should finish successfully", func(ctx ginkgo.SpecContext) {
			ginkgo.By("applying new vulnerability finding that references a package")
			findingInfo := apitypes.FindingInfo{}
			err := findingInfo.FromVulnerabilityFindingInfo(apitypes.VulnerabilityFindingInfo{
				Package: &apitypes.Package{
					Name:    to.Ptr("pkgName"),
					Version: to.Ptr("pkgVersion"),
				},
				Severity:          to.Ptr(apitypes.HIGH),
				VulnerabilityName: to.Ptr("Vulnerability"),
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			_, err = client.PostFinding(ctx, apitypes.Finding{
				FindingInfo: &findingInfo,
				LastSeen:    to.Ptr(time.Now()),
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			reportFailedConfig.objects = append(reportFailedConfig.objects, APIObject{
				"findings", "id ne null",
			})

			ginkgo.By("waiting until at least one findings summary has been updated")
			gomega.Eventually(func() bool {
				findings, err := client.GetFindings(ctx, apitypes.GetFindingsParams{
					Filter: to.Ptr("summary/updatedAt ne null and summary/totalVulnerabilities/totalHighVulnerabilities eq 1"),
					Top:    to.Ptr(1),
					Count:  to.Ptr(true),
				})
				gomega.Expect(skipDBLockedErr(err)).NotTo(gomega.HaveOccurred())
				return *findings.Count > 0
			}, cfg.TestSuiteParams.ScanTimeout, DefaultPeriod).Should(gomega.BeTrue())
		})
	})

	ginkgo.AfterEach(func(ctx ginkgo.SpecContext) {
		if ginkgo.CurrentSpecReport().Failed() {
			reportFailedConfig.startTime = ginkgo.CurrentSpecReport().StartTime
			ReportFailed(ctx, testEnv, client, &reportFailedConfig)
		}
	})
})
