package controllers

import (
	"elevators/core"
	"elevators/models/states"
	"math/rand"
)

// Calls the elevator with queue with the least elements
type Random struct {
	MaxFloors int
}

func (controller *Random) renewCabin(queue states.SimpleQueue) states.SimpleQueue {
	for i := len(queue); i < 100; i++ {
		queue = append(queue, core.Floor(1 + rand.Intn(controller.MaxFloors - 1)))
	}
	return queue
}

func (controller *Random) NewCall(cur_states []core.StateBase, floors []core.Floor, floor core.Floor, panel int) (int, core.StateBase) {
	if controller.MaxFloors <= 1 {
		panic("trying to use Random controller with MaxFloors <= 1")
	}
	least_ind := 0
	for ind, el := range cur_states {
		if len(el.(states.SimpleQueue)) < len(cur_states[least_ind].(states.SimpleQueue)) {
			least_ind = ind
		}
	}
	state_copy := make([]core.Floor, len(cur_states[least_ind].(states.SimpleQueue)))
	copy(state_copy, cur_states[least_ind].(states.SimpleQueue))
	state_copy = controller.renewCabin(state_copy)
	return least_ind, states.SimpleQueue(state_copy)
}

func (controller *Random) CabinFloorCall(state core.StateBase, floor core.Floor) core.StateBase {
	return controller.renewCabin(state.(states.SimpleQueue))
}

func (controller *Random) IntentionToPanelState(core.Floor, core.Floor) int {
	return 0
}
