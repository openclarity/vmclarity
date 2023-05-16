package integration_test

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/cli"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openclarity/vmclarity/integration_test/testenv"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"testing"
)

var (
	testEnv   *testenv.Environment
	log       logr.Logger
	logSyncFn func() error
	client    *backendclient.BackendClient
)

func TestIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Run integration tests")
}

func beforeSuite(ctx context.Context) {
	var err error

	By("creating test environment")
	log, logSyncFn, err = testenv.NewLogger(GinkgoWriter)
	ctx = logr.NewContext(ctx, log)

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

	u, err := testEnv.VMClarityURL()
	Expect(err).NotTo(HaveOccurred())

	client, err = backendclient.Create(fmt.Sprintf("%s://%s/api", u.Scheme, u.Host))
	Expect(err).NotTo(HaveOccurred())
}

var _ = BeforeSuite(beforeSuite)

func afterSuite(ctx context.Context) {
	By("tearing down test environment")
	err := testEnv.Stop(ctx)
	Expect(err).NotTo(HaveOccurred())
	defer func(fn func() error) {
		err := fn()
		if err != nil {
			fmt.Printf("calling sync on logger failed: %v\n", err)
		}
	}(logSyncFn)
}

var _ = AfterSuite(afterSuite)

func areServicesReady(ctx context.Context) bool {
	ready, err := testEnv.Ready(ctx)
	Expect(err).NotTo(HaveOccurred())
	return ready
}
