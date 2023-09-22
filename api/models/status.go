package models

import (
	"fmt"
	"time"
)

type Defaulter interface {
	Default()
}

func NewResourceCleanupStatus(state ResourceCleanupStatusState, m string, s ...string) *ResourceCleanupStatus {
	lastTransitionTime := time.Now()
	var reason string
	if len(s) == 1 {
		reason = s[0]
	}

	return &ResourceCleanupStatus{
		LastTransitionTime: &lastTransitionTime,
		State:              &state,
		Message:            &m,
		Reason:             &reason,
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

func (r *ResourceCleanupStatus) UpdateState(toState ResourceCleanupStatusState) error {
	if *r.State == toState {
		return nil
	}

	transitions, ok := ResourceCleanupStatusStateValidTransitions[*r.State]
	if ok {
		for _, transition := range transitions {
			if transition == toState {
				r.State = &toState
				return nil
			}
		}
	}

	return fmt.Errorf("invalid transition: %s > %s", *r.State, toState)
}

func (r *ResourceCleanupStatus) Default() {
	now := time.Now()
	state := ResourceCleanupStatusStatePending

	if r == nil {
		*r = *NewResourceCleanupStatus(state, "Default ResourceCleanupStatus created.")
	}

	if r.State == nil {
		r.LastTransitionTime = &now
		r.State = &state
	}
}
