package models

import "time"

func NewSbomScan(p *[]Package, s SbomScanState, r SbomScanReason, m *string) *SbomScan {
	return &SbomScan{
		Packages:           p,
		State:              s,
		Reason:             r,
		Message:            m,
		LastTransitionTime: time.Now(),
	}
}
