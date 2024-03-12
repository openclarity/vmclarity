package types

import (
	"context"
	"sync"
)

// TODO(ramizpolic): figure out interface
type JobManager interface {
	Add(ctx context.Context, scanID string, scanTemplate ScanTemplate) error
	Start() error
	Stop() error
	Done() bool
	Wait()
	GetScan(scanID string) (*Scan, error)
	GetResult(scanID string) (*ScanResult, error)
	GetScanner() Scanner
	SetScanner(s Scanner)
}

var _ JobManager = &manager{}

type manager struct {
	mu sync.RWMutex

	scan    Scan
	result  ScanResult
	scanner Scanner
}

func (m *manager) GetInfo() (*ScannerInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Scan(ctx context.Context, input ScanObjectInput, resultCh chan<- Result) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Start() error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) GetResult() (*ScanResult, error) {
	//TODO implement me
	panic("implement me")
}

func NewScanner(template ScanTemplate) ScanManager {
	return &manager{
		mu:      sync.RWMutex{},
		scan:    Scan{},
		result:  ScanResult{},
		scanner: nil,
	}
}

func (m *manager) StartScan(template ScanTemplate) (*Scan, error) {
	switch m.scan.Status.State {
	case ScanStatusStatePending, ScanStatusStateInProgress:
		return nil, ErrScanAlreadyExists
	}

	m.result = ScanResult{
		Annotations: &Annotations{},
		Findings:    []ScanFinding{},
		Summary:     &ScanSummary{},
	}
}

func (m *manager) GetScan() (*Scan, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) StopScan() error {
	//TODO implement me
	panic("implement me")
}

func (m *manager) GetScanResult() (*ScanResult, error) {
	//TODO implement me
	panic("implement me")
}
