package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/integration_test/helpers"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
	"net/http"
	"time"
)

var _ = Describe("Running a default scan (SBOM, vulnerabilities and exploits)", func() {

	Context("which scans a docker container", func() {
		It("should finish successfully", func(ctx SpecContext) {

			By("applying a scan configuration")
			apiScanConfig, err := client.PostScanConfig(ctx, helpers.GetDefaultScanConfig())
			Expect(err).NotTo(HaveOccurred())

			By("waiting until grype server is ready")
			Eventually(func() bool {
				_, err = http.Get("http://localhost:9991")
				return err == nil
			}, time.Second*600).Should(BeTrue())

			By("waiting until trivy server is ready")
			Eventually(func() bool {
				_, err = http.Get("http://localhost:9992")
				return err == nil
			}, time.Second*600).Should(BeTrue())

			By("waiting until exploit db server is ready")
			Eventually(func() bool {
				_, err = http.Get("http://localhost:1326")
				return err == nil
			}, time.Second*600).Should(BeTrue())

			By("updating a scan configuration to run now")
			updateScanConfig := helpers.UpdateScanConfigToStartNow(apiScanConfig)
			err = client.PatchScanConfig(ctx, *apiScanConfig.Id, updateScanConfig)
			Expect(err).NotTo(HaveOccurred())

			By("waiting until scan starts")
			scanParams := models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state ne '%s' and state ne '%s'",
					*apiScanConfig.Id,
					models.ScanStateDone,
					models.ScanStateFailed,
				)),
			}
			var scans *models.Scans
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, scanParams)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*60, time.Second).Should(BeTrue())

			By("waiting until scan state changes to done")
			scanParams = models.GetScansParams{
				Filter: utils.PointerTo(fmt.Sprintf(
					"scanConfig/id eq '%s' and state eq '%s'",
					*apiScanConfig.Id,
					models.ScanStateDone,
				)),
			}
			Eventually(func() bool {
				scans, err = client.GetScans(ctx, scanParams)
				Expect(err).NotTo(HaveOccurred())
				return len(*scans.Items) == 1
			}, time.Second*120, time.Second).Should(BeTrue())
		})
	})
})
