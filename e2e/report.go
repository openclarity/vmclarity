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
	"encoding/json"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/ginkgo/v2/formatter"
	"github.com/onsi/gomega"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/e2e/testenv"
	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

type APIObject struct {
	objectType string
	filter     string
}

type ReportFailedConfig struct {
	// if true, print logs for all services
	allServices bool
	// if not empty, print logs for services in slice
	services []string
	// if not empty, print the last n service logs. if empty, print all.
	serviceLogsTail string
	// if true, print all assets
	allAPIAssets bool
	// if true, print all scan configs
	allAPIScanConfigs bool
	// if true, print all scans
	allAPIScans bool
	// if not empty, print objects in slice
	objects []APIObject
}

// ReportFailed gathers relevant API data and docker service logs for debugging purposes.
func ReportFailed(ctx ginkgo.SpecContext, testEnv *testenv.Environment, client *backendclient.BackendClient, config *ReportFailedConfig) {
	ginkgo.GinkgoWriter.Println("------------------------------")

	DumpAPIData(ctx, client, config)
	DumpServiceLogs(ctx, testEnv, config)

	ginkgo.GinkgoWriter.Println("------------------------------")
}

// nolint:cyclop
// DumpAPIData prints API objects filtered using test parameters (e.g. assets filtered by scope, scan configs filtered by id).
// If filter not provided, no objects are printed.
func DumpAPIData(ctx ginkgo.SpecContext, client *backendclient.BackendClient, config *ReportFailedConfig) {
	ginkgo.GinkgoWriter.Println(formatter.F("{{red}}[FAILED] Report API Data:{{/}}"))

	if config.allAPIAssets {
		config.objects = append(config.objects, APIObject{"asset", ""})
	}

	if config.allAPIScanConfigs {
		config.objects = append(config.objects, APIObject{"scanConfigs", ""})
	}

	if config.allAPIScans {
		config.objects = append(config.objects, APIObject{"scans", ""})
	}

	for _, object := range config.objects {
		switch object.objectType {
		case "asset":
			var params models.GetAssetsParams
			if object.filter == "" {
				params = models.GetAssetsParams{}
			} else {
				params = models.GetAssetsParams{Filter: utils.PointerTo(object.filter)}
			}
			assets, err := client.GetAssets(ctx, params)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			buf, err := json.Marshal(*assets.Items)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			ginkgo.GinkgoWriter.Printf("Asset: %s\n", string(buf))

		case "scanConfig":
			var params models.GetScanConfigsParams
			if object.filter == "" {
				params = models.GetScanConfigsParams{}
			} else {
				params = models.GetScanConfigsParams{Filter: utils.PointerTo(object.filter)}
			}
			scanConfigs, err := client.GetScanConfigs(ctx, params)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			buf, err := json.Marshal(*scanConfigs.Items)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			ginkgo.GinkgoWriter.Printf("Scan Config: %s\n", string(buf))

		case "scan":
			var params models.GetScansParams
			if object.filter == "" {
				params = models.GetScansParams{}
			} else {
				params = models.GetScansParams{Filter: utils.PointerTo(object.filter)}
			}
			scans, err := client.GetScans(ctx, params)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			buf, err := json.Marshal(*scans.Items)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			ginkgo.GinkgoWriter.Printf("Scan: %s\n", string(buf))
		}
	}
}

// DumpServiceLogs prints logs for all services.
func DumpServiceLogs(ctx ginkgo.SpecContext, testEnv *testenv.Environment, config *ReportFailedConfig) {
	ginkgo.GinkgoWriter.Println(formatter.F("{{red}}[FAILED] Report Service Logs:{{/}}"))

	var services []string
	if config.allServices {
		services = testEnv.Services()
	} else {
		services = config.services
	}

	tail := config.serviceLogsTail
	if len(tail) == 0 {
		tail = "all"
	}

	err := testEnv.ServicesLogs(ctx, services, tail, formatter.ColorableStdOut, formatter.ColorableStdErr)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
