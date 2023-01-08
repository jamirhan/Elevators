package states

import (
	"elevators/core"
)


type SimpleQueue []core.Floor

func (state SimpleQueue) FloorReachedResponse(floor core.Floor) (core.Response, core.StateBase) {
	if len(state) == 0 {
		return core.Response{
			Direction: core.Stay,
			Open:      false,
		}, state
	}
	var dir core.Direction
	var open bool

	if floor == state[0] {
		open = true
		state = (state)[1:]
	}

	if len(state) == 0 {
		dir = core.Stay
	} else {
		if (state)[0] > floor {
			dir = core.Up
		}

		if (state)[0] < floor {
			dir = core.Down
		}
	}

	return core.Response{
		Direction: dir,
		Open:      open,
	}, state

}

func DefaultSimpleQueue() SimpleQueue {
	return make(SimpleQueue, 0)
}

func InsertBefore(queue SimpleQueue, index int, el core.Floor) (SimpleQueue) {
	res := queue[:index]
	res = append(res, el)
	return append(res, queue[index:]...)
}

func PushBack(queue SimpleQueue, el core.Floor) SimpleQueue {
	return InsertBefore(queue, len(queue), el)
}

func PushFront(queue SimpleQueue, el core.Floor) SimpleQueue {
	return InsertBefore(queue, 0, el)
}
