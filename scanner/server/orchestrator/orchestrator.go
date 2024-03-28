package orchestrator

import (
	"errors"
	"github.com/openclarity/vmclarity/scanner/server/scanner"
	"github.com/openclarity/vmclarity/scanner/types"
)

var ErrNotRunning = errors.New("scan not running")

type Orchestrator interface {
	Scanner() scanner.Scanner
	Start() error
	Stop() error
	StartScan(scan types.Scan) error
	StopScan(scanID string) error
}
