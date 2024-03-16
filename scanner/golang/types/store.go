package types

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type PreconditionFailedError struct {
	Reason string
}

func (e *PreconditionFailedError) Error() string {
	return "Precondition failed: " + e.Reason
}

type ScanStore interface {
	GetScans(params GetScansParams) (Scans, error)
	GetScan(scanID ScanID) (Scan, error)
	CreateScan(scan Scan) (Scan, error)
	UpdateScan(scanID ScanID, scan Scan) (Scan, error)
}

type ScanResultStore interface {
	GetScanResult(scanID ScanID) (ScanResult, error)
	CreateScanResult(result ScanResult) (ScanResult, error)
	UpdateScanResult(scanID ScanID, result ScanResult) (ScanResult, error)
}

type Store interface {
	ScanStore() ScanStore
	ScanResultStore() ScanResultStore
}
