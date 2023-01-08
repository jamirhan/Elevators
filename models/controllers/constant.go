package controllers

import "elevators/core"


// Always calls the same elevator, no mater how many are present.
// Adding floor to elevator's state, as well as type of elevator's state, is user-defined.
type Constant[T core.StateBase] struct {
	Index int
	NewFloor func(T, core.Floor) T
}

func (controller *Constant[T]) NewCall(states []core.StateBase, floors []core.Floor, floor core.Floor, panel int) (int, core.StateBase) {
	return controller.Index, controller.NewFloor(states[0].(T), floor)
}

func (controller *Constant[T]) CabinFloorCall(state core.StateBase, floor core.Floor) core.StateBase {
	return controller.NewFloor(state.(T), floor)
}

func (controller *Constant[T]) IntentionToPanelState(core.Floor, core.Floor) int {
	return 0
}
