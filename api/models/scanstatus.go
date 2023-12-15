package models

import (
	"fmt"
	"time"
)

var scanStatusStateTransitions = map[ScanStatusState][]ScanStatusState{
	ScanStatusStatePending: {
		ScanStatusStateDiscovered,
		ScanStatusStateAborted,
		ScanStatusStateFailed,
		ScanStatusStateDone,
	},
	ScanStatusStateDiscovered: {
		ScanStatusStateInProgress,
		ScanStatusStateAborted,
	},
	ScanStatusStateInProgress: {
		ScanStatusStateAborted,
		ScanStatusStateFailed,
		ScanStatusStateDone,
	},
	ScanStatusStateAborted: {
		ScanStatusStateFailed,
	},
}

var scanStatusReasonMapping = map[ScanStatusState][]ScanStatusReason{
	ScanStatusStatePending: {
		ScanStatusReasonCreated,
	},
	ScanStatusStateDiscovered: {
		ScanStatusReasonAssetsDiscovered,
	},
	ScanStatusStateInProgress: {
		ScanStatusReasonAssetScansRunning,
	},
	ScanStatusStateAborted: {
		ScanStatusReasonCancellation,
	},
	ScanStatusStateFailed: {
		ScanStatusReasonCancellation,
		ScanStatusReasonAssetScanFailed,
		ScanStatusReasonTimeout,
	},
	ScanStatusStateDone: {
		ScanStatusReasonNothingToScan,
		ScanStatusReasonSuccess,
	},
}

func NewScanStatus(s ScanStatusState, r ScanStatusReason, m *string) *ScanStatus {
	return &ScanStatus{
		State:              s,
		Reason:             r,
		Message:            m,
		LastTransitionTime: time.Now(),
	}
}

func (a *ScanStatus) Equals(aa *ScanStatus) bool {
	if a.Message == nil && aa.Message != nil {
		return false
	}
	if aa.Message == nil && a.Message != nil {
		return false
	}
	if a.Message == nil && aa.Message == nil {
		return a.State == aa.State && a.Reason == aa.Reason
	}

	return a.State == aa.State && a.Reason == aa.Reason && *a.Message == *aa.Message
}

func (a *ScanStatus) isValidStatusTransition(aa *ScanStatus) error {
	transitions, ok := scanStatusStateTransitions[a.State]
	if ok {
		for _, transition := range transitions {
			if transition == aa.State {
				return nil
			}
		}
	}

	return fmt.Errorf("invalid transition: from=%s to=%s", a.State, aa.State)
}

func (a *ScanStatus) isValidReason(aa *ScanStatus) error {
	reasons, ok := scanStatusReasonMapping[aa.State]
	if ok {
		for _, reason := range reasons {
			if reason == aa.Reason {
				return nil
			}
		}
	}

	return fmt.Errorf("invalid reason for state: state=%s reason=%s", a.State, a.Reason)
}

func (a *ScanStatus) IsValidTransition(aa *ScanStatus) error {
	if a.Equals(aa) {
		return nil
	}

	if err := a.isValidStatusTransition(aa); err != nil {
		return err
	}
	if err := a.isValidReason(aa); err != nil {
		return err
	}

	return nil
}
