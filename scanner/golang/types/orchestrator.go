package types

import "errors"

var ErrNotRunning = errors.New("scan not running")

type Orchestrator interface {
	Scanner() Scanner
	Store() Store
	Start() error
	Stop() error
	StartScan(scanID string) error
	StopScan(scanID string) error
}
