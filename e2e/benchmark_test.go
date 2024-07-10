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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/openclarity/vmclarity/scanner"
	scannercommon "github.com/openclarity/vmclarity/scanner/common"
	"github.com/openclarity/vmclarity/scanner/families"
	exploits "github.com/openclarity/vmclarity/scanner/families/exploits/types"
	infofinder "github.com/openclarity/vmclarity/scanner/families/infofinder/types"
	malware "github.com/openclarity/vmclarity/scanner/families/malware/types"
	misconfigurations "github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	"github.com/openclarity/vmclarity/scanner/families/plugins/runner/config"
	plugins "github.com/openclarity/vmclarity/scanner/families/plugins/types"
	rootkits "github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	secrets "github.com/openclarity/vmclarity/scanner/families/secrets/types"
	vulnerabilities "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/types"
	"github.com/openclarity/vmclarity/utils/fsutils/containerrootfs"
)

const (
	rootfsPath          = "/tmp/bench-test"
	testImageSourcePath = "./testdata/alpine-3.18.2.tar"
	alpineImage         = "alpine:3.18.2"
	markDownFilePath    = "/tmp/scanner-benchmark.md"
	markdownHeader      = "# 🚀 Benchmark results\n\n"
)

type BenchmarkNotifier struct {
	mu            sync.Mutex
	started       map[families.FamilyType]time.Time
	finished      map[families.FamilyType]time.Time
	findingsCount map[families.FamilyType]int
}

func (n *BenchmarkNotifier) FamilyStarted(_ context.Context, famType families.FamilyType) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.started[famType] = time.Now().UTC()

	return nil
}

// nolint:cyclop
func (n *BenchmarkNotifier) FamilyFinished(_ context.Context, res families.FamilyResult) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.finished[res.FamilyType] = time.Now().UTC()

	switch res.FamilyType {
	case families.SBOM:
		// if res.Result.(*sbom.Result) != nil {
		// 	familyResult := res.Result.(*sbom.Result).SBOM.Vulnerabilities
		// 	n.findingsCount[res.FamilyType] = len(*familyResult)
		// } else {
		// 	n.findingsCount[res.FamilyType] = 0
		// }

	case families.Vulnerabilities:
		if res.Result.(*vulnerabilities.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*vulnerabilities.Result).MergedVulnerabilitiesByKey // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}

	case families.Secrets:
		if res.Result.(*secrets.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*secrets.Result).Findings // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}

	case families.Exploits:
		if res.Result.(*exploits.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*exploits.Result).Exploits // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}

	case families.Misconfiguration:
		if res.Result.(*misconfigurations.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*misconfigurations.Result).Misconfigurations // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}

	case families.Rootkits:
		if res.Result.(*rootkits.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*rootkits.Result).Rootkits // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}

	case families.Malware:
		if res.Result.(*malware.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*malware.Result).Malwares // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}

	case families.InfoFinder:
		if res.Result.(*infofinder.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*infofinder.Result).Infos // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}

	case families.Plugins:
		if res.Result.(*plugins.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*plugins.Result).Findings // nolint:forcetypeassert
			n.findingsCount[res.FamilyType] = len(familyResult)
		} else {
			n.findingsCount[res.FamilyType] = 0
		}
	}

	return nil
}

var _ = ginkgo.Describe("Running a Benchmark test", func() {
	ginkgo.Context("which scans an alpine image", func() {
		ginkgo.It("should finish successfully", func(ctx ginkgo.SpecContext) {
			fmt.Println("Starting benchmark test")
			scannerConfig := &scanner.Config{
				// SBOM: sbom.Config{
				// 	Enabled: true,
				// 	AnalyzersList: []string{
				// 		"syft",
				// 		"trivy",
				// 		"windows",
				// 	},
				// },
				// Vulnerabilities: vulnerabilities.Config{
				// 	Enabled: true,
				// 	ScannersList: []string{
				// 		"grype",
				// 		"trivy",
				// 	},
				// },
				// Secrets: secrets.Config{
				// 	Enabled: true,
				// 	ScannersList: []string{
				// 		"gitleaks",
				// 	},
				// },
				// Exploits: exploits.Config{
				// 	Enabled: true,
				// 	ScannersList: []string{
				// 		"exploitdb",
				// 	},
				// },
				// Misconfiguration: misconfigurations.Config{
				// 	Enabled: true,
				// 	ScannersList: []string{
				// 		"lynis",
				// 		"cisdocker",
				// 		"fake",
				// 	},
				// },
				// Rootkits: rootkits.Config{
				// 	Enabled: true,
				// 	ScannersList: []string{
				// 		"chkrootkit",
				// 	},
				// },
				// Malware: malware.Config{
				// 	Enabled: true,
				// 	ScannersList: []string{
				// 		"clam",
				// 		"yara",
				// 	},
				// },
				// InfoFinder: infofinder.Config{
				// 	Enabled: true,
				// 	ScannersList: []string{
				// 		"sshTopology",
				// 	},
				// },
				Plugins: plugins.Config{
					Enabled:      true,
					ScannersList: []string{scannerPluginNameKICS},
					ScannersConfig: plugins.ScannersConfig{
						scannerPluginNameKICS: config.Config{
							Name:          scannerPluginNameKICS,
							ImageName:     cfg.TestEnvConfig.Images.PluginKics,
							InputDir:      "",
							ScannerConfig: "",
						},
					},
				},
			}

			image, cleanup, err := containerrootfs.GetImageWithCleanup(ctx, testImageSourcePath)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer cleanup()

			err = containerrootfs.ToDirectory(ctx, image, rootfsPath)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			scannerConfig.AddInputs(scannercommon.ROOTFS, []string{rootfsPath})
			scannerConfig.AddInputs(scannercommon.IMAGE, []string{alpineImage})

			input, err := filepath.Abs("./testdata")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			scannerConfig.AddInputs(scannercommon.ROOTFS, []string{input})

			notifier := &BenchmarkNotifier{
				started:       make(map[families.FamilyType]time.Time),
				finished:      make(map[families.FamilyType]time.Time),
				findingsCount: make(map[families.FamilyType]int),
			}

			errs := scanner.New(scannerConfig).Run(ctx, notifier)
			gomega.Expect(errs).To(gomega.BeEmpty())

			mdTable, err := notifier.GenerateMarkdownTable()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = writeMarkdownTableToFile(mdTable)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
})

func (n *BenchmarkNotifier) GenerateMarkdownTable() (string, error) {
	rows := [][]string{}

	earliestStartTime := time.Now().UTC()
	totalFindings := 0
	for famType, startTime := range n.started {
		if startTime.Before(earliestStartTime) {
			earliestStartTime = startTime
		}

		row := []string{
			string(famType),
			n.started[famType].Format(time.DateTime),
			n.finished[famType].Format(time.DateTime),
			strconv.Itoa(n.findingsCount[famType]),
			n.finished[famType].Sub(startTime).Round(time.Second).String(),
		}
		rows = append(rows, row)

		totalFindings += n.findingsCount[famType]
	}

	latestFinishTime := time.Time{}
	for _, endTime := range n.finished {
		if endTime.After(latestFinishTime) {
			latestFinishTime = endTime
		}
	}

	mdTable, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("Family/Scanner", "Start time", "End time", "Findings", "Total time").
		Format(rows)
	if err != nil {
		return "", fmt.Errorf("failed to format markdown table: %w", err)
	}

	footer := fmt.Sprintf(
		"\n\nFull scan summary\nTotal time: %s\nTotal findings: %d\n",
		latestFinishTime.Sub(earliestStartTime).Round(time.Second).String(),
		totalFindings,
	)

	return mdTable + footer, nil
}

func writeMarkdownTableToFile(mdTable string) error {
	err := os.WriteFile(markDownFilePath, []byte(markdownHeader+mdTable), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}

	return nil
}
