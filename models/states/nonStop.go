package states

import "elevators/core"

type NonStop struct {
	curDirection core.Direction
	MaxFloors    int
}

func (state NonStop) FloorReachedResponse(floor core.Floor) (core.Response, core.StateBase) {
	if state.MaxFloors <= 1 {
		panic("trying to use Random controller with MaxFloors <= 1")
	}
	if floor == 0 {
		return core.Response{
				Direction: core.Up,
				Open:      true,
			}, NonStop{
				curDirection: core.Up,
				MaxFloors:    state.MaxFloors,
			}
	}
	if state.curDirection == core.Down {
		return core.Response{
				Direction: core.Down,
				Open:      true,
			}, NonStop{
				curDirection: core.Down,
				MaxFloors:    state.MaxFloors,
			}
	}
	if floor == core.Floor(state.MaxFloors) {
		return core.Response{
				Direction: core.Down,
				Open:      true,
			}, NonStop{
				curDirection: core.Down,
				MaxFloors:    state.MaxFloors,
			}
	}
	if state.curDirection != core.Up {
		panic("incorrect internal state")
	}
	return core.Response{
			Direction: core.Up,
			Open:      true,
		}, NonStop{
			curDirection: core.Up,
			MaxFloors:    state.MaxFloors,
		}
}
