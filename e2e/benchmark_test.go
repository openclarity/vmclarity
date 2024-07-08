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
	misconfigurations "github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	plugins "github.com/openclarity/vmclarity/scanner/families/plugins/types"
	rootkits "github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	sbom "github.com/openclarity/vmclarity/scanner/families/sbom/types"
	secrets "github.com/openclarity/vmclarity/scanner/families/secrets/types"
	vulnerabilities "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/types"
	"github.com/openclarity/vmclarity/utils/fsutils/containerrootfs"
)

const (
	rootfsPath              = "/tmp/bench-test"
	testImage               = "./testdata/alpine-3.18.2.tar"
	defaultMarkDownFilePath = "scanner-benchmark.md"
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

func (n *BenchmarkNotifier) FamilyFinished(_ context.Context, res families.FamilyResult) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.finished[res.FamilyType] = time.Now().UTC()

	switch res.FamilyType {
	case families.SBOM:
		// familyResult := res.Result.(*sbom.Result).SBOM.Vulnerabilities
		// n.findingsCount[res.FamilyType] = len(familyResult)

	case families.Vulnerabilities:
		familyResult := res.Result.(*vulnerabilities.Result).MergedVulnerabilitiesByKey
		n.findingsCount[res.FamilyType] = len(familyResult)

	case families.Secrets:
		familyResult := res.Result.(*secrets.Result).Findings
		n.findingsCount[res.FamilyType] = len(familyResult)

	case families.Exploits:
		familyResult := res.Result.(*exploits.Result).Exploits
		n.findingsCount[res.FamilyType] = len(familyResult)

	case families.Misconfiguration:
		familyResult := res.Result.(*misconfigurations.Result).Misconfigurations
		n.findingsCount[res.FamilyType] = len(familyResult)

	case families.Rootkits:
		familyResult := res.Result.(*rootkits.Result).Rootkits
		n.findingsCount[res.FamilyType] = len(familyResult)

	case families.Malware:
		familyResult := res.Result.(*malware.Result).Malwares
		n.findingsCount[res.FamilyType] = len(familyResult)

	case families.InfoFinder:
		familyResult := res.Result.(*infofinder.Result).Infos
		n.findingsCount[res.FamilyType] = len(familyResult)

	case families.Plugins:
		familyResult := res.Result.(*plugins.Result).Findings
		n.findingsCount[res.FamilyType] = len(familyResult)
	}

	return nil
}

func (n *BenchmarkNotifier) GenerateMarkdownTable() (string, error) {
	rows := [][]string{}
	for famType, startTime := range n.started {
		row := []string{
			string(famType),
			n.started[famType].String(),
			n.finished[famType].String(),
			strconv.Itoa(n.findingsCount[famType]),
			n.finished[famType].Sub(startTime).String(),
		}
		rows = append(rows, row)
	}

	mdTable, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("Family", "Start time", "End time", "Findings", "Total time").
		Format(rows)
	if err != nil {
		return "", fmt.Errorf("failed to format markdown table: %w", err)
	}

	return mdTable, nil
}

func Test_Benchmark(t *testing.T) {
	g := gomega.NewWithT(t)

	fmt.Println("Starting benchmark test")
	scannerConfig := &scanner.Config{
		SBOM: sbom.Config{
			Enabled: true,
			AnalyzersList: []string{
				"syft",
				"trivy",
				"windows",
			},
		},
		Vulnerabilities: vulnerabilities.Config{
			Enabled: true,
		},
		Secrets: secrets.Config{
			Enabled: true,
		},
		Exploits: exploits.Config{
			Enabled: true,
		},
		Misconfiguration: misconfigurations.Config{
			Enabled: true,
		},
		Rootkits: rootkits.Config{
			Enabled: true,
		},
		Malware: malware.Config{
			Enabled: true,
		},
		InfoFinder: infofinder.Config{
			Enabled: true,
		},
		Plugins: plugins.Config{
			Enabled: true,
		},
	}

	ctx := context.Background()
	image, cleanup, err := containerrootfs.GetImageWithCleanup(ctx, testImage)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer cleanup()

	err = containerrootfs.ToDirectory(ctx, image, rootfsPath)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	scannerConfig.AddInputs(scannercommon.ROOTFS, []string{rootfsPath})
	scannerConfig.AddInputs(scannercommon.IMAGE, []string{testImage})

	notifier := &BenchmarkNotifier{
		started:       make(map[families.FamilyType]time.Time),
		finished:      make(map[families.FamilyType]time.Time),
		findingsCount: make(map[families.FamilyType]int),
	}

	errs := scanner.New(scannerConfig).Run(ctx, notifier)
	g.Expect(errs).To(gomega.BeEmpty())

	mdTable, err := notifier.GenerateMarkdownTable()
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = writeMarkdownFile(mdTable)
	g.Expect(err).NotTo(gomega.HaveOccurred())

}

func writeMarkdownFile(mdTable string) error {
	err := os.WriteFile(defaultMarkDownFilePath, []byte(mdTable), 0644)
	if err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}

	return nil
}
