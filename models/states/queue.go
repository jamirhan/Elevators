package states

import (
	"elevators/core"
)

type SimpleQueue[T any] []core.Floor

func (state *SimpleQueue[T]) GetResponse(floor core.Floor) core.Response {
	if len(*state) == 0 {
		return core.Response{
			Direction: core.Stay,
			Open:      false,
		}
	}
	var dir core.Direction
	var open bool

	if floor == (*state)[0] {
		open = true
		*state = (*state)[1:]
	}

	if len(*state) == 0 {
		dir = core.Stay
	} else {
		if (*state)[0] > floor {
			dir = core.Up
		}

		if (*state)[0] < floor {
			dir = core.Down
		}
	}

	return core.Response{
		Direction: dir,
		Open:      open,
	}

}

func (state *SimpleQueue[T]) NewFloor(floor core.Floor) {
	*state = append(*state, floor)
}

func (state *SimpleQueue[T]) NewCall(floor core.Floor, panel T) {
	*state = append(*state, floor)
}

func DefaultSimpleQueue[T any]() *SimpleQueue[T] {
	var queue SimpleQueue[T]
	return &queue
}