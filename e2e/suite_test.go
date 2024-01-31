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
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/e2e/config"
	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
	"github.com/openclarity/vmclarity/pkg/shared/uibackendclient"
	"github.com/openclarity/vmclarity/testenv"
	"github.com/openclarity/vmclarity/testenv/types"
	"github.com/openclarity/vmclarity/utils/log"
)

var (
	testEnv types.Environment
	cfg     *config.Config

	client   *backendclient.BackendClient
	uiClient *uibackendclient.UIBackendClient
)

func TestEndToEnd(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Run end-to-end tests")
}

func beforeSuite(ctx context.Context) {
	var err error

	ginkgo.By("initializing test environment")
	log.InitLogger(logrus.DebugLevel.String(), ginkgo.GinkgoWriter)
	logger := logrus.WithContext(ctx)
	ctx = log.SetLoggerForContext(ctx, logger)

	// Get end-to-end test suite configuration
	cfg, err = config.NewConfig()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Create new testenv from configuration
	testEnv, err = testenv.New(&cfg.TestEnvConfig, testenv.WithContext(ctx), testenv.WithLogger(logger))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	if !cfg.ReuseEnv {
		setupCtx, cancel := context.WithTimeout(ctx, cfg.EnvSetupTimeout)
		defer cancel()

		ginkgo.By("setup test environment")
		err = testEnv.SetUp(setupCtx)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	} else {
		ginkgo.By("re-using test environment")
	}

	ginkgo.By("waiting for services to become ready")
	gomega.Eventually(func() bool {
		ready, err := testEnv.ServicesReady(ctx)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		return ready
	}, time.Second*5).Should(gomega.BeTrue())

	endpoints, err := testEnv.Endpoints(ctx)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	client, err = backendclient.Create(endpoints.API.String())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	uiClient, err = uibackendclient.Create(endpoints.UIBackend.String())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

var _ = ginkgo.BeforeSuite(beforeSuite)

func afterSuite(ctx context.Context) {
	if !cfg.ReuseEnv {
		ginkgo.By("tearing down test environment")
		err := testEnv.TearDown(ctx)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
}

var _ = ginkgo.AfterSuite(afterSuite)
