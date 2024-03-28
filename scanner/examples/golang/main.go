package main

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/client"
	"github.com/openclarity/vmclarity/scanner/server/orchestrator"
	"golang_example/scanner"
	"os"
	"os/signal"
	"syscall"
)

// TODO: make this cleaner

const (
	VMCLARITY_SCANNER_SERVER = "http://0.0.0.0:8765"
)

// labels that this scanner uses to be able to take scans available on the
// scanner server all labels must be present in the scan attributes/metadata in
// order to be picked up by this scanner (AND operator).
var watchPendingScanLabels = map[string]string{
	"scanner/name":   "cisdocker",
	"scanner/family": "misconfigurations",
}

func main() {
	err := mainE()
	if err != nil {
		fmt.Println("error", err)
	}
}

func mainE() error {
	ctx := context.Background()

	// Create client
	scannerClient, err := client.NewClient(VMCLARITY_SCANNER_SERVER)
	if err != nil {
		return err
	}

	manager, err := orchestrator.NewOrchestrator(
		ctx,
		&scanner.Scanner{},
		scannerClient,
		watchPendingScanLabels,
	)
	if err != nil {
		return err
	}

	_ = manager.Start()
	defer manager.Stop()

	// Wait for shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sig
	fmt.Printf("Received a termination signal: %v\n", s)

	return nil
}
