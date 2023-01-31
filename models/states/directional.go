package states;

import (
	"elevators/core"
)


type Directional struct {
	direction core.Direction
	less map[core.Floor]bool
	more map[core.Floor]bool
	lastFloor core.Floor
}


func (state Directional) FloorReachedResponse(floor core.Floor) (core.Response, core.StateBase) {
	state.lastFloor = floor;
	if _, ok := state.more[floor]; ok {
		delete(state.more, floor)
		dir := core.Up;
		if len(state.more) == 0 {
			dir = core.Down
		}
		if len(state.less) == 0 {
			dir = core.Stay
		}
		state.direction = dir
		return core.Response{
			Direction: dir,
			Open: true,
		}, state
	} else if _, ok := state.less[floor]; ok {
		delete(state.less, floor)
		dir := core.Down;
		if len(state.less) == 0 {
			dir = core.Up
		}
		if len(state.more) == 0 {
			dir = core.Stay
		}
		state.direction = dir
		return core.Response{
			Direction: dir,
			Open: true,
		}, state
	} else {
		return core.Response{
			Direction: state.direction,
			Open: false,
		}, state
	}
}


func (state *Directional) Add(floor core.Floor) {
	if state.lastFloor >= floor {
		state.less[floor] = true
	} else {
		state.more[floor] = true
	}
}