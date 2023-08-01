package integration_test

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/cli"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/integration_test/testenv"
	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
	"os"
	"strconv"
	"testing"
	"time"
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
		[]string{"../installation/docker/docker-compose.yaml"},
		cli.WithName("vmclarity"),
		cli.WithWorkingDirectory("../installation/docker"),
		cli.WithResolvedPaths(true),
		cli.WithProfiles([]string{""}),
	)
	Expect(err).NotTo(HaveOccurred())

	_, err = os.Create(os.ExpandEnv("$HOME/vmclarity.db"))
	Expect(err).NotTo(HaveOccurred())

	err = cli.WithOsEnv(opts)
	Expect(err).NotTo(HaveOccurred())

	var reuseEnv bool
	if reuseEnv, _ = strconv.ParseBool(os.Getenv("USE_EXISTING")); reuseEnv {
		log.V(-1).Info("reusing existing environment...", "use_existing", reuseEnv)
	}

	testEnv, err = testenv.New(opts, reuseEnv)
	Expect(err).NotTo(HaveOccurred())

	By("starting test environment")
	err = testEnv.Start(ctx)
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() bool {
		ready, err := testEnv.ServicesReady(ctx)
		Expect(err).NotTo(HaveOccurred())
		return ready
	}, time.Second*5).Should(BeTrue())

	u, err := testEnv.VMClarityURL()
	Expect(err).NotTo(HaveOccurred())

	client, err = backendclient.Create(fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, u.Path))
	Expect(err).NotTo(HaveOccurred())

	// todo(adam.tagscherer): create a proper readyz endpoint for the api
	By("waiting until VMClarity API is ready")
	Eventually(func() bool {
		_, err = client.GetScanConfigs(ctx, models.GetScanConfigsParams{})
		return err == nil
	}, time.Second*5).Should(BeTrue())
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
