package models

import (
	"fmt"
	"time"
)

var resourceCleanupStatusStateTransitions = map[ResourceCleanupStatusState][]ResourceCleanupStatusState{
	ResourceCleanupStatusStatePending: {
		ResourceCleanupStatusStateSkipped,
		ResourceCleanupStatusStateFailed,
		ResourceCleanupStatusStateDone,
	},
}

var resourceCleanupStatusReasonMapping = map[ResourceCleanupStatusState][]ResourceCleanupStatusReason{
	ResourceCleanupStatusStatePending: {
		ResourceCleanupStatusReasonAssetScanCreated,
	},
	ResourceCleanupStatusStateSkipped: {
		ResourceCleanupStatusReasonDeletePolicy,
	},
	ResourceCleanupStatusStateFailed: {
		ResourceCleanupStatusReasonProviderError,
		ResourceCleanupStatusReasonInternalError,
	},
	ResourceCleanupStatusStateDone: {
		ResourceCleanupStatusReasonSuccess,
	},
}

func NewResourceCleanupStatus(s ResourceCleanupStatusState, r ResourceCleanupStatusReason, m *string) *ResourceCleanupStatus {
	return &ResourceCleanupStatus{
		State:              s,
		Reason:             r,
		Message:            m,
		LastTransitionTime: time.Now(),
	}
}

func NewResourceCleanupStatusWithDefaults() *ResourceCleanupStatus {
	return NewResourceCleanupStatus(ResourceCleanupStatusStatePending, ResourceCleanupStatusReasonAssetScanCreated, nil)
}

func (rs *ResourceCleanupStatus) Equals(r *ResourceCleanupStatus) bool {
	return rs.State == r.State && rs.Reason == r.Reason && *rs.Message == *r.Message
}

func (rs *ResourceCleanupStatus) IsValidTransition(r *ResourceCleanupStatus) error {
	if rs.Equals(r) {
		return nil
	}

	transitions, ok := resourceCleanupStatusStateTransitions[rs.State]
	var isValid bool
	if ok {
		for _, transition := range transitions {
			if transition == r.State {
				isValid = true
				break
			}
		}
	}
	if !ok || !isValid {
		return fmt.Errorf("invalid transition: from=%s to=%s", rs.State, r.State)
	}

	reasons, ok := resourceCleanupStatusReasonMapping[r.State]
	if ok {
		for _, reason := range reasons {
			if reason == r.Reason {
				isValid = true
				break
			}
		}
	}
	if !ok || !isValid {
		return fmt.Errorf("invalid reason for state: state=%s reason=%s", r.State, r.Reason)
	}

	return nil
}
