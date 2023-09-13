package models

import "fmt"

type ResourceCleanupStateMachine map[ResourceCleanupStatusState][]ResourceCleanupStatusState

var ResourceCleanupStatusStateValidTransitions = ResourceCleanupStateMachine{
	ResourceCleanupStatusStatePending: {
		ResourceCleanupStatusStateSkipped,
		ResourceCleanupStatusStateFailed,
		ResourceCleanupStatusStateDone,
	},
}

func (rcs *ResourceCleanupStatus) UpdateState(toState ResourceCleanupStatusState) error {
	transitions := ResourceCleanupStatusStateValidTransitions[*rcs.State]
	for _, transition := range transitions {
		if transition == toState {
			rcs.State = &toState
			return nil
		}
	}

	return fmt.Errorf("invalid transition, from: %s to: %s", *rcs.State, toState)
}
