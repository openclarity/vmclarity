package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/openclarity/vmclarity/scanner/client"
	"github.com/openclarity/vmclarity/scanner/pkg/scanner"
	"github.com/openclarity/vmclarity/scanner/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"sync"
	"time"
)

// TODO: add options struct to make everything more configurable

const (
	maxScanJobs        = 100              // total number of scans that can run at the same time
	maxQueueSize       = 1000             // total number of scans that can be queued
	scanSchedulePeriod = 10 * time.Second // period to check for pending tasks and schedule them
)

var ErrScanQueueFull = fmt.Errorf("scan queue is full and cannot accept new tasks yet")

type manager struct {
	sync.RWMutex
	ctx        context.Context
	scanner    scanner.Scanner
	scanLabels map[string]string
	client     *client.Client
	workerPool *pond.WorkerPool
	scanCancel map[string]func()
	cancel     func()
}

// NewOrchestrator creates a new Orchestrator. Pass scanMetadataSelector to only
// watch pending scans that contain given labels.
func NewOrchestrator(ctx context.Context,
	scanner scanner.Scanner,
	client *client.Client,
	watchPendingScanLabels map[string]string,
) (Orchestrator, error) {
	log.SetLevel(log.DebugLevel)

	// Validate
	if len(watchPendingScanLabels) == 0 {
		return nil, fmt.Errorf("orchestrator should not watch everything, pass scan metadata labels to watch")
	}

	// Create orchestrator context
	ctx, cancel := context.WithCancel(ctx)

	// Create a worker pool
	pool := pond.New(maxScanJobs, maxQueueSize, pond.Context(ctx), pond.PanicHandler(func(err interface{}) {
		log.Errorf("worker recovered from panic, reason: %v", err)
	}))

	// Return initialized orchestrator
	return &manager{
		ctx:        ctx,
		scanner:    scanner,
		scanLabels: watchPendingScanLabels,
		client:     client,
		workerPool: pool,
		scanCancel: make(map[string]func()),
		cancel:     cancel,
	}, nil
}

func (m *manager) Scanner() scanner.Scanner { return m.scanner }

func (m *manager) Start() error {
	scanLabelsSelector := types.MetaSelectorsFromMap(m.scanLabels)
	go func() {
		ticker := time.NewTicker(scanSchedulePeriod)
		defer ticker.Stop()

		scanPendingState := types.ScanState(types.ScanStatusStatePending)
		for {
			select {
			case <-ticker.C:
				log.WithField("scan-labels", m.scanLabels).Infof("checking for pending tasks...")

				// Get pending tasks
				scans, err := respTo[types.Scans](m.client.GetScans(m.ctx, &types.GetScansParams{
					State:         &scanPendingState,
					MetaSelectors: scanLabelsSelector,
				}))
				if err != nil {
					log.Errorf("failed to get pending scans: %v", err)
					continue
				}

				// Start pending tasks
				for _, scan := range scans.Items {
					scanID := *scan.Id
					if err = m.StartScan(scan); err != nil {
						log.Errorf("orchestrator failed to schedule scan %s: %v", scanID, err)
						continue
					}
					log.Infof("successfully queued scan %s", scanID)
				}

			case <-m.ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (m *manager) Stop() error {
	// TODO: add timeout here and run in goroutine; wait for ctx done

	m.cancel()
	m.workerPool.Stop()

	m.Lock()
	for _, cancelScan := range m.scanCancel {
		if cancelScan != nil {
			cancelScan()
		}
	}
	m.scanCancel = map[string]func(){}
	m.Unlock()

	return nil
}

// StartScan adds scan job to the worker queue. Workers will pick up the scan job
// when they become available. If the worker queue is full, this method returns
// an error notifying the caller that the scan could not be queued. Workers
// execute multiple scan jobs in parallel. Each scan job performs its own
// sub-scans in parallel too, but waits for all sub-scans to complete before
// marking the scan job as completed and releasing the worker. This will not wait
// for the queue to become free and will drop instantly.
func (m *manager) StartScan(scan types.Scan) error {
	queued := m.workerPool.TrySubmit(func() {
		scanID := *scan.Id
		log.Infof("Scan %s started...", scanID)

		// Create sub-scan processor
		scanCtx, scanCancel := context.WithCancel(m.ctx)                                // with cancel
		if scan.InProgressTimeoutSeconds != nil && *scan.InProgressTimeoutSeconds > 0 { // with timeout
			duration := time.Duration(*scan.InProgressTimeoutSeconds) * time.Second
			scanCtx, scanCancel = context.WithTimeout(scanCtx, duration)
		}
		scanJob, scanCtx := errgroup.WithContext(scanCtx)
		defer scanCancel()

		// Send handshake event or return
		{
			scannerInfo := m.scanner.GetInfo(scanCtx)
			scanEvent := &types.ScanEvent_EventInfo{}
			err := scanEvent.FromScannerHandshakeEventInfo(types.ScannerHandshakeEventInfo{
				Annotations: scannerInfo.Annotations,
				EventType:   "Handshake", // TODO: expose proper models that does this automatically
				Name:        scannerInfo.Name,
			})
			if err != nil {
				log.Errorf("Could not send handshake event to the scanner server")
				return
			}
			_, err = m.client.SubmitScanEvent(m.ctx, scanID, types.SubmitScanEventJSONRequestBody{
				EventInfo: *scanEvent,
			})
			if err != nil {
				log.Errorf("Could not send handshake event to the scanner server")
				return
			}
		}

		// Add running scan cancellation function to the internal map
		m.Lock()
		m.scanCancel[scanID] = scanCancel
		m.Unlock()

		// Remove scan from internal map when completed/exited
		defer func() {
			m.Lock()
			delete(m.scanCancel, scanID)
			m.Unlock()
		}()

		// Use intermediary output to avoid calling db on each sub-scan completion. In
		// case of issues, we should probably monitor this as it might cause high RAM
		// spikes due to size of the worker pool. By default, we drop all sub-scan
		// results when any of them fails.
		inputScanCh := make(chan types.ScanFinding)

		// Run all sub-scans in parallel within a given worker. The worker will become
		// available for new scan jobs when this one is completed.
		var inputs []types.ScanInput
		if scan.Inputs != nil {
			inputs = *scan.Inputs
		}
		for _, input := range inputs {
			input := input
			scanJob.Go(func() error {
				var err error

				// Handle panics
				defer func() {
					if panicErr := recover(); panicErr != nil {
						err = fmt.Errorf("runtime panic, reason: %v", panicErr)
					}
				}()

				// Run scan
				results, err := m.Scanner().Scan(scanCtx, scanID, input)
				if err != nil {
					return err
				}

				// Send results via channel to track progress
				for _, result := range results {
					inputScanCh <- result
				}

				return nil
			})
		}

		// Start scan in goroutine to allow watching the results
		errCh := make(chan error, 1)
		go func() {
			errCh <- scanJob.Wait()
			close(inputScanCh)
		}()

		// Watch for scan results and update stats
		results := make([]types.ScanFinding, 0, len(inputs))
		for subScan := range inputScanCh {
			results = append(results, subScan)

			// TODO: pass analytics data to the scanner server as event here
		}

		// Extract details from scan result
		state := types.ScannerHeartbeatCompleted
		stateMsg := "scan completed successfully"
		summary := types.ScanSummary{
			JobsDone: len(results),
		}
		if jobErr := <-errCh; jobErr != nil {
			switch {
			case errors.Is(jobErr, context.Canceled):
				state = types.ScannerHeartbeatCancelled
				stateMsg = "scan aborted, reason: cancelled"
			case errors.Is(jobErr, context.DeadlineExceeded):
				state = types.ScannerHeartbeatCancelled
				stateMsg = "scan aborted, reason: timed-out"
			default:
				state = types.ScannerHeartbeatErrored
				stateMsg = fmt.Sprintf("scan failed, reason: %v", jobErr)
			}
			summary.JobsFailed = 1
			log.Errorf("Scan %s errored while processing, reason: %v", scanID, jobErr)
		}

		// Submit findings event
		{
			findingsEvent := &types.ScanEvent_EventInfo{}
			err := findingsEvent.FromScannerFindingsEventInfo(types.ScannerFindingsEventInfo{
				Findings: results,
			})
			if err != nil {
				log.Errorf("Could not send findings event to the scanner server")
				return
			}
			_, err = m.client.SubmitScanEvent(m.ctx, scanID, types.SubmitScanEventJSONRequestBody{
				EventInfo: *findingsEvent,
			})
			if err != nil {
				log.Errorf("Could not send findings event to the scanner server")
				return
			}
		}

		// Submit final heartbeat
		{
			heartbeatEvent := &types.ScanEvent_EventInfo{}
			err := heartbeatEvent.FromScannerHeartbeatEventInfo(types.ScannerHeartbeatEventInfo{
				Message: &stateMsg,
				State:   state,
				Summary: &summary,
			})
			if err != nil {
				log.Errorf("Could not send final heartbeat event to the scanner server")
				return
			}
			_, err = m.client.SubmitScanEvent(m.ctx, scanID, types.SubmitScanEventJSONRequestBody{
				EventInfo: *heartbeatEvent,
			})
			if err != nil {
				log.Errorf("Could not send final heartbeat event to the scanner server")
				return
			}
		}

		log.Infof("Scan %s successfully processed, discovered %d findings", scanID, len(results))
	})
	if !queued {
		return ErrScanQueueFull
	}

	return nil
}

func (m *manager) StopScan(scanID string) error {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.scanCancel[scanID]; !ok {
		return ErrNotRunning
	}

	m.Lock()
	if cancelFn := m.scanCancel[scanID]; cancelFn != nil {
		cancelFn()
	}
	delete(m.scanCancel, scanID)
	m.Unlock()

	return nil
}

func toPtr[T any](t T) *T {
	return &t
}

func respTo[T any](r *http.Response, err error) (T, error) {
	var target T

	// check original error
	if err != nil {
		return target, fmt.Errorf("failed while performing request: %w", err)
	}

	// decode
	defer r.Body.Close()
	bytes, err := io.ReadAll(r.Body)

	err = json.Unmarshal(bytes, &target) //.Decode(&target)
	if err != nil {
		return target, fmt.Errorf("failed to parse result to %T: %w", target, err)
	}

	return target, nil
}
