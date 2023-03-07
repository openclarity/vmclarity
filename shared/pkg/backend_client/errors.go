package backend_client

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

type TargetConflictError struct {
	ConflictingTarget *models.Target
	Message           string
}

func (t TargetConflictError) Error() string {
	return fmt.Sprintf("Conflicting Target Found with ID %s: %s", *t.ConflictingTarget.Id, t.Message)
}

type ScanConflictError struct {
	ConflictingScan *models.Scan
	Message         string
}

func (t ScanConflictError) Error() string {
	return fmt.Sprintf("Conflicting Scan Found with ID %s: %s", *t.ConflictingScan.Id, t.Message)
}

type ScanResultConflictError struct {
	ConflictingScanResult *models.TargetScanResult
	Message               string
}

func (t ScanResultConflictError) Error() string {
	return fmt.Sprintf("Conflicting Scan Result Found with ID %s: %s", *t.ConflictingScanResult.Id, t.Message)
}
