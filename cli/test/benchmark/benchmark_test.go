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

package benchmark

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/onsi/gomega"

	"github.com/openclarity/vmclarity/scanner"
	scannercommon "github.com/openclarity/vmclarity/scanner/common"
	"github.com/openclarity/vmclarity/scanner/families"
	exploits "github.com/openclarity/vmclarity/scanner/families/exploits/types"
	infofinder "github.com/openclarity/vmclarity/scanner/families/infofinder/types"
	malware "github.com/openclarity/vmclarity/scanner/families/malware/types"
	yaraconfig "github.com/openclarity/vmclarity/scanner/families/malware/yara/config"
	misconfigurations "github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	plugins "github.com/openclarity/vmclarity/scanner/families/plugins/types"
	rootkits "github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	sbom "github.com/openclarity/vmclarity/scanner/families/sbom/types"
	secrets "github.com/openclarity/vmclarity/scanner/families/secrets/types"
	grypeconfig "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/grype/config"
	trivyconfig "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/trivy/config"
	vulnerabilities "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/types"
	"github.com/openclarity/vmclarity/utils/fsutils/containerrootfs"
)

const (
	rootfsPath          = "/tmp/bench-test"
	testImageSourcePath = "./testdata/alpine-3.18.2.tar"
	markDownFilePath    = "/tmp/scanner-benchmark.md"
	tableHeader         = "# 🚀 Benchmark results"
)

type BenchmarkInfo struct {
	Name          string
	StartTime     time.Time
	EndTime       time.Time
	TotalTime     time.Duration
	FindingsCount string
}

type BenchmarkNotifier struct {
	mu       sync.Mutex
	families map[families.FamilyType]BenchmarkInfo
	scanners map[families.FamilyType][]BenchmarkInfo
}

func (n *BenchmarkNotifier) FamilyStarted(_ context.Context, famType families.FamilyType) error {
	return nil
}

// nolint:cyclop
func (n *BenchmarkNotifier) FamilyFinished(_ context.Context, res families.FamilyResult) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	var findingsCount int
	var metadata families.ScanMetadata
	switch res.FamilyType {
	case families.SBOM:
		if res.Result.(*sbom.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*sbom.Result) // nolint:forcetypeassert
			if familyResult.SBOM.Components != nil {
				findingsCount = len(*familyResult.SBOM.Components)
			} else {
				findingsCount = 0
			}
			metadata = familyResult.Metadata
		}

	case families.Vulnerabilities:
		if res.Result.(*vulnerabilities.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*vulnerabilities.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.MergedVulnerabilitiesByKey)
			metadata = familyResult.Metadata
		}

	case families.Secrets:
		if res.Result.(*secrets.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*secrets.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.Findings)
			metadata = familyResult.Metadata
		}

	case families.Exploits:
		if res.Result.(*exploits.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*exploits.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.Exploits)
			metadata = familyResult.Metadata
		}

	case families.Misconfiguration:
		if res.Result.(*misconfigurations.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*misconfigurations.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.Misconfigurations)
			metadata = familyResult.Metadata
		}

	case families.Rootkits:
		if res.Result.(*rootkits.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*rootkits.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.Rootkits)
			metadata = familyResult.Metadata
		}

	case families.Malware:
		if res.Result.(*malware.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*malware.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.Malwares)
			metadata = familyResult.Metadata
		}

	case families.InfoFinder:
		if res.Result.(*infofinder.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*infofinder.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.Infos)
			metadata = familyResult.Metadata
		}

	case families.Plugins:
		if res.Result.(*plugins.Result) != nil { // nolint:forcetypeassert
			familyResult := res.Result.(*plugins.Result) // nolint:forcetypeassert
			findingsCount = len(familyResult.Findings)
			metadata = familyResult.Metadata
		}
	}

	n.families[res.FamilyType] = BenchmarkInfo{
		Name:          fmt.Sprintf("%s/*", res.FamilyType),
		StartTime:     metadata.StartTime,
		EndTime:       metadata.EndTime,
		TotalTime:     metadata.EndTime.Sub(metadata.StartTime),
		FindingsCount: strconv.Itoa(findingsCount),
	}

	if len(metadata.Inputs) == 0 {
		return nil
	}

	n.scanners[res.FamilyType] = append(n.scanners[res.FamilyType], BenchmarkInfo{
		Name:          fmt.Sprintf("%s/%s", res.FamilyType, metadata.Inputs[0].ScannerName),
		StartTime:     metadata.StartTime,
		EndTime:       metadata.EndTime,
		TotalTime:     metadata.EndTime.Sub(metadata.StartTime),
		FindingsCount: "-",
	})

	return nil
}

func Test_Benchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping benchmark test in short mode")
	}

	t.Run("Benchmark test", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		ctx := context.Background()

		scannerConfig := &scanner.Config{
			SBOM: sbom.Config{
				Enabled: true,
				AnalyzersList: []string{
					"syft",
					"trivy",
					"windows",
					"gomod",
				},
			},
			Vulnerabilities: vulnerabilities.Config{
				Enabled: true,
				ScannersList: []string{
					"grype",
					"trivy",
				},
				ScannersConfig: vulnerabilities.ScannersConfig{
					Grype: grypeconfig.Config{
						Mode: grypeconfig.ModeLocal,
						Local: grypeconfig.LocalGrypeConfig{
							UpdateDB:      true,
							UpdateTimeout: 5 * time.Minute,
						},
					},
					Trivy: trivyconfig.Config{
						Timeout: 300,
					},
				},
			},
			Secrets: secrets.Config{
				Enabled: true,
				ScannersList: []string{
					"gitleaks",
				},
			},
			// TODO: enable exploits once it has been added to base tools
			// Exploits: exploits.Config{
			// 	Enabled: true,
			// 	Inputs: []scannercommon.ScanInput{
			// 		{
			// 			Input:     "CVE-2006-2896,CVE-2007-2007",
			// 			InputType: scannercommon.CSV,
			// 		},
			// 	},
			// 	ScannersList: []string{
			// 		"exploitdb",
			// 	},
			// 	ScannersConfig: exploits.ScannersConfig{
			// 		ExploitDB: exploitsconfig.Config{
			// 			BaseURL: "http://127.0.0.1:1326",
			// 		},
			// 	},
			// },
			Misconfiguration: misconfigurations.Config{
				Enabled: true,
				ScannersList: []string{
					"cisdocker",
					"lynis",
					"fake",
				},
			},
			Rootkits: rootkits.Config{
				Enabled: true,
				ScannersList: []string{
					"chkrootkit",
				},
			},
			Malware: malware.Config{
				Enabled: true,
				ScannersList: []string{
					"clam",
					"yara",
				},
				ScannersConfig: malware.ScannersConfig{
					Yara: yaraconfig.Config{
						CompiledRuleURL: "https://raw.githubusercontent.com/Yara-Rules/rules/master/malware/APT_APT1.yar",
					},
				},
			},
			InfoFinder: infofinder.Config{
				Enabled: true,
				ScannersList: []string{
					"sshTopology",
				},
			},
			// TODO: enable plugins once the issues with the runner are fixed
			// Plugins: plugins.Config{
			// 	Enabled: true,
			// 	ScannersList: []string{
			// 		"kics",
			// 	},
			// 	Inputs: []scannercommon.ScanInput{
			// 		{
			// 			Input:     "../../../e2e/testdata",
			// 			InputType: scannercommon.ROOTFS,
			// 		},
			// 	},
			// 	ScannersConfig: plugins.ScannersConfig{
			// 		"kics": pluginsconfig.Config{
			// 			Name:      "kics",
			// 			ImageName: "ghcr.io/openclarity/vmclarity-plugin-kics:latest",
			// 		},
			// 	},
			// },
		}

		image, cleanup, err := containerrootfs.GetImageWithCleanup(ctx, testImageSourcePath)
		g.Expect(err).NotTo(gomega.HaveOccurred())
		defer cleanup()

		err = containerrootfs.ToDirectory(ctx, image, rootfsPath)
		g.Expect(err).NotTo(gomega.HaveOccurred())

		scannerConfig.AddInputs(scannercommon.ROOTFS, []string{rootfsPath})

		notifier := &BenchmarkNotifier{
			families: make(map[families.FamilyType]BenchmarkInfo),
			scanners: make(map[families.FamilyType][]BenchmarkInfo),
		}

		errs := scanner.New(scannerConfig).Run(ctx, notifier)
		g.Expect(errs).To(gomega.BeEmpty())

		mdTable, err := notifier.GenerateMarkdownTable()
		g.Expect(err).NotTo(gomega.HaveOccurred())

		err = writeMarkdownTableToFile(mdTable)
		g.Expect(err).NotTo(gomega.HaveOccurred())
	})
}

func (n *BenchmarkNotifier) GenerateMarkdownTable() (string, error) {
	rows := [][]string{}
	var earliestStartTime, latestEndTime time.Time
	totalFindingsCount := 0
	for famType, benchmarkInfo := range n.families {
		if earliestStartTime.IsZero() || benchmarkInfo.StartTime.Before(earliestStartTime) {
			earliestStartTime = benchmarkInfo.StartTime
		}

		if latestEndTime.IsZero() || benchmarkInfo.EndTime.After(latestEndTime) {
			latestEndTime = benchmarkInfo.EndTime
		}

		findingsCount, err := strconv.Atoi(benchmarkInfo.FindingsCount)
		if err != nil {
			return "", fmt.Errorf("failed to convert findings count to integer: %w", err)
		}
		totalFindingsCount += findingsCount

		rows = append(rows, []string{
			benchmarkInfo.Name,
			benchmarkInfo.StartTime.Format(time.DateTime),
			benchmarkInfo.EndTime.Format(time.DateTime),
			benchmarkInfo.FindingsCount,
			benchmarkInfo.TotalTime.Round(time.Second).String(),
		})

		for _, scannerInfo := range n.scanners[famType] {
			rows = append(rows, []string{
				scannerInfo.Name,
				scannerInfo.StartTime.Format(time.DateTime),
				scannerInfo.EndTime.Format(time.DateTime),
				scannerInfo.FindingsCount,
				scannerInfo.TotalTime.Round(time.Second).String(),
			})
		}
	}

	// append a separator row and the summary
	rows = append(rows, []string{"", "", "", "", ""}, []string{
		"_Scan summary_",
		fmt.Sprintf("_%s_", earliestStartTime.Format(time.DateTime)),
		fmt.Sprintf("_%s_", latestEndTime.Format(time.DateTime)),
		fmt.Sprintf("_%d_", totalFindingsCount),
		fmt.Sprintf("_%s_", latestEndTime.Sub(earliestStartTime).Round(time.Second).String()),
	})

	tableBody, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("Family/Scanner", "Start time", "End time", "Findings", "Total time").
		Format(rows)
	if err != nil {
		return "", fmt.Errorf("failed to format markdown table body: %w", err)
	}

	return tableHeader + "\n\n" + tableBody, nil
}

func writeMarkdownTableToFile(mdTable string) error {
	err := os.WriteFile(markDownFilePath, []byte(mdTable), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}

	return nil
}
