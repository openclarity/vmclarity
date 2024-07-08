//go:build benchmark
// +build benchmark

package benchmark

import (
	"fmt"
	"testing"
)

func TestBenchmark(t *testing.T) {
	fmt.Println("benchmark test")
}

// import (
// 	"context"
// 	"sync"
// 	"testing"
// 	"time"
// )

// type Notifier struct {
// 	mu            sync.Mutex
// 	started       map[families.FamilyType]time.Time
// 	finished      map[families.FamilyType]time.Time
// 	findingsCount map[families.FamilyType]int
// }

// func (n *Notifier) FamilyStarted(_ context.Context, famType families.FamilyType) error {
// 	n.started[famType.FamilyType] = time.Now()

// 	return nil
// }

// func (n *Notifier) FamilyFinished(_ context.Context, famType families.FamilyResult) error {
// 	n.finished[famType.FamilyType] = time.Since(n.started[famType.FamilyType])

// 	switch famType.FamilyType {
// 	case families.Exploits:
// 		result, _ := res.Result.(*types.Result)

// 		n.findingCount[res.FamilyType] = len(result.Exploits)
// 	}

// 	return nil
// }

// func (n *Notifier) GenerateMarkdownFile(outputMarkdownFile string) error {

// 	return nil
// }

// func TestBenchmark(t *testing.T) {
// 	notifier := &Notifier{
// 		started:       make(map[families.FamilyType]time.Time),
// 		finished:      make(map[families.FamilyType]time.Time),
// 		findingsCount: make(map[families.FamilyType]int),
// 	}
// 	// 1. Create config without inputs, just enable every scanner

// 	// Need to set enabled and ScannerList/AnalyzerList

// 	scannerConfig := &Config{}

// 	// 2. Create ROOTFS under /tmp/bench-test from an image.

// 	// Use containerrootfs.GetImageWithCleanup and containerrootfs.ToDirectory

// 	// 3. Add inputs to the scanner config

// 	// We want to test ROOTFS and IMAGE

// 	scannerConfig.AddInputs(common.ROOTFS, "rootfs-from-the-testimage")

// 	scannerConfig.AddInputs(common.IMAGE, "testimage")

// 	// 4. Run the scan

// 	errs := New(scannerConfig).Run(ctx, notifier)

// 	gomega.Expect(errs).To(gomega.BeEmpty())

// 	// 5. Save to output

// 	notifier.GenerateMarkdownFile()

// 	// 6. Push the generated file to artifacts under "scanner-benchmark.md"

// 	// 7. Write summary to GH action run with https://github.blog/2022-05-09-supercharging-github-actions-with-job-summaries/

// 	// TABLE FORMAT:

// 	// Family | Start time (time: HH:MM:SS) | End time (time: HH:MM:SS) | Findings (COUNT) | Total time (duration: 5m20s)
// }
