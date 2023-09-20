package models

import (
	"fmt"
	"time"
)

type Defaulter interface {
	Default()
}

func NewResourceCleanupStatus(lastTransitionTime time.Time, message, reason string, state ResourceCleanupStatusState) *ResourceCleanupStatus {
	return &ResourceCleanupStatus{
		LastTransitionTime: &lastTransitionTime,
		Message:            &message,
		Reason:             &reason,
		State:              &state,
	}
}

type ResourceCleanupStateMachine map[ResourceCleanupStatusState][]ResourceCleanupStatusState

var ResourceCleanupStatusStateValidTransitions = ResourceCleanupStateMachine{
	ResourceCleanupStatusStatePending: {
		ResourceCleanupStatusStateSkipped,
		ResourceCleanupStatusStateFailed,
		ResourceCleanupStatusStateDone,
	},
}

func (rcs *ResourceCleanupStatus) UpdateState(toState ResourceCleanupStatusState) error {
	if *rcs.State == toState {
		return nil
	}

	transitions, ok := ResourceCleanupStatusStateValidTransitions[*rcs.State]
	if ok {
		for _, transition := range transitions {
			if transition == toState {
				rcs.State = &toState
				return nil
			}
		}
	}

	return fmt.Errorf("invalid transition: %s > %s", *rcs.State, toState)
}

func (rcs *ResourceCleanupStatus) Default() {
	now := time.Now()
	state := ResourceCleanupStatusStatePending

	if rcs == nil {
		*rcs = ResourceCleanupStatus{
			LastTransitionTime: &now,
			Message:            nil,
			Reason:             nil,
			State:              &state,
		}
	}

	if rcs.State == nil {
		rcs.LastTransitionTime = &now
		rcs.State = &state
	}
}
