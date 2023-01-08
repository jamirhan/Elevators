package core

import (
	"sync"
)

type state struct {
	state        StateBase
	stateMut     sync.Mutex
	curFloor     Floor
	curDirection Direction
}

func (state *state) initialize(base StateBase) {
	state.state = base
	state.curFloor = 1
	state.curDirection = Stay
}

type elevator struct {
	stateChanged    chan bool
	floorReached    chan ElevatorArrivedT
	curState        *state
	interFloorTicks int
	tickChan        chan bool
	elInd           int
}

func startElevator(stateChanged chan bool, floorReached chan ElevatorArrivedT, tickChan chan bool, initState *state, ind int) {
	lift := elevator {
		stateChanged: stateChanged,
		tickChan:     tickChan,
		floorReached: floorReached,
		curState:     initState,
		elInd:        ind,
	}
	go lift.elevatorRoutine()
}

func (lift *elevator) elevatorRoutine() {
	for {
		lift.elevatorIteration()
	}
}

func (lift *elevator) elevatorIteration() {
	select {
	case <-lift.tickChan:

		lift.curState.stateMut.Lock()
		defer lift.curState.stateMut.Unlock()

		if lift.curState.curDirection == Stay {
			return
		}
		if lift.interFloorTicks == TicksBetweenFloors {
			lift.interFloorTicks = 0
			if lift.curState.curDirection == Up {
				lift.curState.curFloor++
			} else {
				lift.curState.curFloor--
			}
			resp, new_state := lift.curState.state.FloorReachedResponse(lift.curState.curFloor)
			lift.curState.state = new_state
			lift.curState.curDirection = resp.Direction
			if resp.Open {
				arrived := ElevatorArrivedT{Floor: lift.curState.curFloor, ElevatorInd: lift.elInd}
				lift.floorReached <- arrived
			}
		} else {
			lift.interFloorTicks++
		}
	case <-lift.stateChanged:
		lift.curState.stateMut.Lock()
		defer lift.curState.stateMut.Unlock()
		if lift.curState.curDirection == Stay {
			resp, new_state := lift.curState.state.FloorReachedResponse(lift.curState.curFloor)
			lift.curState.state = new_state
			lift.curState.curDirection = resp.Direction
			if resp.Open {
				arrived := ElevatorArrivedT{Floor: lift.curState.curFloor, ElevatorInd: lift.elInd}
				lift.floorReached <- arrived
			}
		}
	}
}
