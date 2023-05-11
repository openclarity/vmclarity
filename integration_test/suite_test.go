package integration_test

import (
	"context"
	"github.com/compose-spec/compose-go/cli"
	"github.com/openclarity/vmclarity/integration_test/testenv"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	testEnv *testenv.Environment
	client  *backendclient.BackendClient
)

func TestIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Run integration tests")
}

func beforeSuite(ctx context.Context) {
	var err error

	By("creating test environment")

	opts, err := cli.NewProjectOptions(
		[]string{"../deploy/docker-compose.yml"},
		cli.WithName("vm-clarity"),
	)
	Expect(err).NotTo(HaveOccurred())

	testEnv, err = testenv.New(opts)
	Expect(err).NotTo(HaveOccurred())

	By("starting test environment")
	err = testEnv.Start(ctx)
	Expect(err).NotTo(HaveOccurred())

	Eventually(areServicesReady(ctx), 15).Should(BeTrue())
}

var _ = BeforeSuite(beforeSuite)

func afterSuite(ctx context.Context) {
	By("tearing down test environment")
	err := testEnv.Stop(ctx)
	Expect(err).NotTo(HaveOccurred())
}

var _ = AfterSuite(afterSuite)

func areServicesReady(ctx context.Context) bool {
	ready, err := testEnv.Ready(ctx)
	Expect(err).NotTo(HaveOccurred())
	return ready
}
