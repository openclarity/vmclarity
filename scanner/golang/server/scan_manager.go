package server

import (
	"github.com/openclarity/vmclarity/scanner/types"
)

type ScanManager interface {
	Start() error
	Stop() error
	StartScan(scanTemplate types.ScanTemplate) (*types.Scan, error)
	StopScan(scanID string) error
	ScanDone(scanID string) bool
	GetScan(scanID string) (*types.Scan, error)
	GetScanResult(scanID string) (*types.ScanResult, error)
	Scanner() types.Scanner
}

var _ ScanManager = &manager{}

type manager struct{}

func (m manager) Start() error {
	//TODO implement me
	panic("implement me")
}

func (m manager) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (m manager) StartScan(scanTemplate types.ScanTemplate) (*types.Scan, error) {
	//TODO implement me
	panic("implement me")
}

func (m manager) StopScan(scanID string) error {
	//TODO implement me
	panic("implement me")
}

func (m manager) ScanDone(scanID string) bool {
	//TODO implement me
	panic("implement me")
}

func (m manager) GetScan(scanID string) (*types.Scan, error) {
	//TODO implement me
	panic("implement me")
}

func (m manager) GetScanResult(scanID string) (*types.ScanResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m manager) Scanner() types.Scanner {
	//TODO implement me
	panic("implement me")
}
