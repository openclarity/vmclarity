package orchestrator

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
	"time"
)

type manager struct {
	scanner types.Scanner
	store   types.Store

	totalJobs atomic.Uint64
	jobs      map[string]chan struct{}
	stopCh    chan struct{}
	stopOnce  sync.Once
}

func NewOrchestrator(scanner types.Scanner, store types.Store) (types.ScanManager, error) {
	return &manager{
		scanner: scanner,
		store:   store,
	}, nil
}

func (m *manager) Scanner() types.Scanner {
	return m.scanner
}

func (m *manager) Store() types.Store {
	return m.store
}

func (m *manager) StartScan(scanID string) error {
	// Get scan
	scan, err := m.Store().ScanStore().GetScan(scanID)
	if err != nil {
		return fmt.Errorf("cannot start scan %s: %v", scanID, err)
	}

	// Check scan state
	switch state := scan.Status.State; state {
	case types.ScanStatusStatePending:
	default:
		return fmt.Errorf("cannot start scan %s due to invalid state %s", scanID, state)
	}

	// Create processor
	m.totalJobs.Add(1)
	procGroup, procCtx := errgroup.WithContext(context.Background())

	for _, input := range scan.Inputs {
		input := input
		procGroup.Go(func() error {
			result, err := m.Scanner().Scan(procCtx, scanID, input)
			if err != nil {
				failErr := fmt.Errorf("scan %s failed while running: %v", scanID, err)
				scan, getErr := m.Store().ScanStore().GetScan(scanID)

				if getErr != nil {
					failErr = fmt.Errorf("%w: failed to fetch scan", scan)
				}

			}
		})
	}

	return nil
}

func (m *manager) StopScan(scanID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) ScanDone(scanID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Start() error {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		maxItems := 50
		pending := types.ScanState(types.ScanStatusStatePending)

		for {
			select {
			case <-ticker.C:
				log.Infof("checking for pending tasks...")

				// Get pending tasks
				scans, err := m.Store().ScanStore().GetScans(types.GetScansParams{
					PageSize: &maxItems,
					State:    &pending,
				})
				if err != nil {
					log.Errorf("failed to get pending scans: %v", err)
					continue
				}

				log.Infof("scheduling pending tasks...")

				// Start pending tasks
				for _, scan := range scans.Items {
					scanID := *scan.Id
					if err = m.StartScan(scanID); err != nil {
						log.Errorf("failed to start scan %s: %v", scanID, err)
						continue
					}
					log.Infof("scan %s successfully scheduled", scanID)
				}

				log.Infof("waiting for next tick to schedule...")

			case <-m.stopCh:
				return
			}
		}
	}()
	return nil
}

func (m *manager) Stop() error {
	m.stopOnce.Do(func() {
		close(m.stopCh)
	})
	return nil
}

func (m *manager) onScanSuccess(scanID string, result []types.ScanFinding) error {
	scan, err := m.Store().ScanStore().GetScan(scanID)
	if err != nil {
		return fmt.Errorf("cannot fetch scan: %w", err)
	}

	msg := "scan finished successfully"
	scan.Status = &types.ScanStatus{
		LastTransitionTime: time.Now(),
		Message:            &msg,
		State:              types.ScanStatusStateDone,
	}
	scan.

}
