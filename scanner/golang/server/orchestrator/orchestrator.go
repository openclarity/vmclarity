package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/openclarity/vmclarity/scanner/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

const (
	maxScanJobs  = 100  // total number of scans that can run at the same time
	maxQueueSize = 1000 // total number of scans that can be queued
)

type manager struct {
	mu sync.RWMutex

	scanner    types.Scanner
	store      types.Store
	workerPool *pond.WorkerPool
	scanCancel map[string]func()
	stopOnce   sync.Once
	stopCh     chan struct{}
}

func NewOrchestrator(scanner types.Scanner, store types.Store) (types.Orchestrator, error) {
	// Create a worker pool
	pool := pond.New(maxScanJobs, maxQueueSize, pond.PanicHandler(func(err interface{}) {
		log.Errorf("worker recovered from panic, reason: %v", err)
	}))

	// Return initialized orchestrator
	return &manager{
		scanner:    scanner,
		store:      store,
		workerPool: pool,
		scanCancel: make(map[string]func()),
	}, nil
}

func (m *manager) Scanner() types.Scanner {
	return m.scanner
}

func (m *manager) Store() types.Store {
	return m.store
}

func (m *manager) StartScan(scanID string) error {
	// Create reusable validator as it needs to be checked in the worker
	validator := func() (types.Scan, error) {
		// Get scan
		scan, err := m.Store().Scans().Get(scanID)
		if err != nil {
			return types.Scan{}, fmt.Errorf("cannot fetch scan %s from store: %v", scanID, err)
		}

		// Check scan state
		switch state := scan.Status.State; state {
		case types.ScanStatusStatePending:
		default:
			return types.Scan{}, fmt.Errorf("cannot queue scan %s due to state %s", scanID, state)
		}

		return scan, nil
	}

	// Do validation before queueing the task as it might already be running
	if _, validateErr := validator(); validateErr != nil {
		return validateErr
	}

	// Add scan job to the worker queue. Workers will pick up the scan job when they
	// become available. If the worker queue is full, this method returns an error
	// notifying the caller that the scan could not be queued. Workers execute
	// multiple scan jobs in parallel. Each scan job performs its own sub-scans in
	// parallel too, but waits for all sub-scans to complete before marking the scan
	// job as completed and releasing the worker. This will not wait for the queue to
	// become free and will drop instantly.
	queued := m.workerPool.TrySubmit(func() {
		log.Infof("Scan %s started...", scanID)

		// Revalidate to ensure that no other workers picked this job yet.
		// Safe to return as the job will be re-queued.
		scan, validateErr := validator()
		if validateErr != nil {
			log.Errorf("Skipped performing scan on %s, reason: %v", scanID, validateErr)
			return
		}

		// Create sub-scan processor
		scanCtx := context.Background()                             // base ctx
		scanCtx, scanCancel := context.WithCancel(scanCtx)          // with cancel
		if scan.TimeoutSeconds != nil && *scan.TimeoutSeconds > 0 { // with timeout
			duration := time.Duration(*scan.TimeoutSeconds) * time.Second
			scanCtx, scanCancel = context.WithTimeout(scanCtx, duration)
		}
		scanJob, scanCtx := errgroup.WithContext(scanCtx)
		defer scanCancel()

		// Add running scan cancellation function to the internal map
		m.mu.Lock()
		m.scanCancel[scanID] = scanCancel
		m.mu.Unlock()

		// Remove scan from internal map when completed/exited
		defer func() {
			m.mu.Lock()
			delete(m.scanCancel, scanID)
			m.mu.Unlock()
		}()

		// Set scan state to in progress. It's okay to return here as the job will be
		// re-queued by the orchestrator.
		scan.Status = &types.ScanStatus{
			LastTransitionTime: time.Now(),
			Message:            toPtr("scan started"),
			State:              types.ScanStatusStateInProgress,
		}
		if err := m.updateScan(scanID, &scan); err != nil {
			log.Errorf("Worker failed to set scan %s in progress, reason: %v", scanID, err)
			return
		}

		// Use intermediary output to avoid calling db on each sub-scan completion. In
		// case of issues, we should probably monitor this as it might cause high RAM
		// spikes due to size of the worker pool. By default, we drop all sub-scan
		// results when any of them fails.
		subScanCh := make(chan types.ScanFinding)

		// Run all sub-scans in parallel within a given worker. The worker will become
		// available for new scan jobs when this one is completed.
		for _, input := range scan.Inputs {
			input := input
			scanJob.Go(func() error {
				var err error
				defer func() { // handle panic
					if panicErr := recover(); panicErr != nil {
						err = fmt.Errorf("runtime panic, reason: %v", panicErr)
					}
				}()

				results, err := m.Scanner().Scan(scanCtx, scanID, input)
				if err != nil {
					return err
				}

				// Send results via channel to track progress
				for _, result := range results {
					subScanCh <- result
				}

				return nil
			})
		}

		// Start scan in goroutine to allow watching the results
		jobErrCh := make(chan error, 1)
		go func() {
			jobErrCh <- scanJob.Wait()
			close(subScanCh)
		}()

		// Watch for scan results and update stats
		jobResults := make([]types.ScanFinding, len(scan.Inputs))
		startTime := time.Now()
		{
			updateStep := 1 + (len(scan.Inputs) / 5) // avoid 0
			for subScan := range subScanCh {
				jobResults = append(jobResults, subScan)

				// Update scan completion once every 1/5th of inputs
				// processed to not overwork the DB
				if len(jobResults)%updateStep == 1 {
					scan.JobsCompleted = toPtr(len(jobResults))
					scan.JobsLeftToRun = toPtr(len(scan.Inputs) - len(jobResults))
					_ = m.updateScan(scanID, &scan)
				}
			}
		}
		endTime := time.Now()

		// Extract details from scan result
		state := types.ScanStatusStateDone
		stateMsg := "scan completed successfully"
		if jobErr := <-jobErrCh; jobErr != nil {
			switch {
			case errors.Is(jobErr, context.Canceled):
				state = types.ScanStatusStateAborted
				stateMsg = "scan aborted, reason: cancelled"
			case errors.Is(jobErr, context.DeadlineExceeded):
				state = types.ScanStatusStateAborted
				stateMsg = "scan aborted, reason: timed-out"
			default:
				state = types.ScanStatusStateFailed
				stateMsg = fmt.Sprintf("scan failed, reason: %v", jobErr)
			}
			log.Errorf("Scan %s errored while processing, reason: %v", scanID, jobErr)
		}

		// Update scan data from the result
		if err := m.updateScan(scanID, &types.Scan{
			EndTime:       &endTime,
			Id:            &scanID,
			Inputs:        scan.Inputs,
			JobsCompleted: toPtr(len(jobResults)),
			JobsLeftToRun: toPtr(len(scan.Inputs) - len(jobResults)),
			StartTime:     &startTime,
			Status: &types.ScanStatus{
				LastTransitionTime: time.Now(),
				Message:            &stateMsg,
				State:              state,
			},
			TimeoutSeconds: scan.TimeoutSeconds,
		}); err != nil {
			log.Errorf("Scan %s completed but its state could not be updated: %v", scanID, err)
			return
		}

		// Create all scan findings
		for idx := range jobResults {
			jobResults[idx].ScanID = &scanID
		}
		if err := m.Store().ScanFindings().CreateMany(jobResults...); err != nil {
			log.Errorf("Scan %s completed but its findings could not be created: %v", scanID, err)
			return
		}

		log.Infof("Scan %s successfully processed, discovered %d findings", scanID, len(jobResults))
	})
	if !queued {
		return fmt.Errorf("workers cannot accept any new tasks")
	}

	return nil
}

func (m *manager) StopScan(scanID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.scanCancel[scanID]; !ok {
		return types.ErrNotRunning
	}

	m.mu.Lock()
	if cancelFn := m.scanCancel[scanID]; cancelFn != nil {
		cancelFn()
	}
	delete(m.scanCancel, scanID)
	m.mu.Unlock()

	return nil
}

func (m *manager) Start() error {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		scanPendingState := types.ScanState(types.ScanStatusStatePending)
		for {
			select {
			case <-ticker.C:
				log.Infof("checking for pending tasks...")

				// Get pending tasks
				scans, err := m.Store().Scans().GetAll(types.GetScansRequest{
					State: &scanPendingState,
				})
				if err != nil {
					log.Errorf("failed to get pending scans: %v", err)
					continue
				}

				// Start pending tasks
				for _, scan := range scans {
					scanID := *scan.Id
					if err = m.StartScan(scanID); err != nil {
						log.Errorf("orchestrator failed to schedule scan %s: %v", scanID, err)
						continue
					}
					log.Infof("successfully queued scan %s", scanID)
				}

			case <-m.stopCh:
				return
			}
		}
	}()
	return nil
}

func (m *manager) Stop() error {
	m.stopOnce.Do(func() {
		m.workerPool.Stop()
		close(m.stopCh)

		m.mu.Lock()
		for _, cancelScan := range m.scanCancel {
			if cancelScan != nil {
				cancelScan()
			}
		}
		m.scanCancel = map[string]func(){}
		m.mu.Unlock()
	})
	return nil
}

func (m *manager) updateScan(scanID string, scan *types.Scan) error {
	scan.Id = &scanID // set scan ID before update
	newScan, err := m.store.Scans().Update(scanID, *scan)
	if err != nil {
		return fmt.Errorf("failed to update scan %s: %w", scanID, err)
	}
	*scan = newScan
	return nil
}

func toPtr[T any](t T) *T {
	return &t
}
