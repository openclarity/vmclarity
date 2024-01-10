package models

import (
	"time"
)

func NewAssetScanEstimationStatus(s AssetScanEstimationStatusState, r AssetScanEstimationStatusReason, m *string) *AssetScanEstimationStatus {
	return &AssetScanEstimationStatus{
		State:              s,
		Reason:             r,
		Message:            m,
		LastTransitionTime: time.Now(),
	}
}
