package controllers

import (
	"elevators/core"
	"elevators/models/states"
)

// Calls the elevator with queue with the least elements
type NonStop struct {
}

func (controller *NonStop) NewCall(cur_states []core.StateBase, floors []core.Floor, floor core.Floor, panel int) (int, core.StateBase) {
	least_ind := 0
	for ind, el := range cur_states {
		if len(el.(states.SimpleQueue)) < len(cur_states[least_ind].(states.SimpleQueue)) {
			least_ind = ind
		}
	}
	return least_ind, states.PushBack(cur_states[least_ind].(states.SimpleQueue), floor)
}

func (controller *NonStop) CabinFloorCall(state core.StateBase, floor core.Floor) core.StateBase {
	return states.PushBack(state.(states.SimpleQueue), floor)
}

func (controller *NonStop) IntentionToPanelState(core.Floor, core.Floor) int {
	return 0
}
