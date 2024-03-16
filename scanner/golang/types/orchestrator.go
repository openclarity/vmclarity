package types

type ScanManager interface {
	Start() error
	Stop() error
	StartScan(scanID string) error
	StopScan(scanID string) error
	ScanDone(scanID string) (bool, error)

	Scanner() Scanner
	Store() Store
}
