package core

import (
	"sync"
)

type state[T any] struct {
	state        StateBase[T]
	stateMut     sync.Mutex
	curFloor     Floor
	curDirection Direction
}

func (state *state[T]) initialize(base StateBase[T]) {
	state.state = base
	state.curFloor = 1
	state.curDirection = Stay
}

type elevator[T any] struct {
	stateChanged    chan bool
	floorReached    chan ElevatorArrivedT
	curState        *state[T]
	interFloorTicks int
	tickChan        chan bool
	elInd           int
}

func startElevator[T any](stateChanged chan bool, floorReached chan ElevatorArrivedT, tickChan chan bool, initState *state[T], ind int) {
	lift := elevator[T]{
		stateChanged: stateChanged,
		tickChan:     tickChan,
		floorReached: floorReached,
		curState:     initState,
		elInd:        ind,
	}
	go lift.elevatorRoutine()
}

func (lift *elevator[T]) elevatorRoutine() {
	for {
		lift.elevatorIteration()
	}
}

func (lift *elevator[T]) elevatorIteration() {
	select {
	case <-lift.tickChan:

		lift.curState.stateMut.Lock()
		defer lift.curState.stateMut.Unlock()

		if lift.curState.curDirection == Stay {
			return
		}
		if lift.interFloorTicks == ticksBetweenFloors {
			lift.interFloorTicks = 0
			if lift.curState.curDirection == Up {
				lift.curState.curFloor++
			} else {
				lift.curState.curFloor--
			}
			resp := lift.curState.state.GetResponse(lift.curState.curFloor)
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
			resp := lift.curState.state.GetResponse(lift.curState.curFloor)
			lift.curState.curDirection = resp.Direction
			if resp.Open {
				arrived := ElevatorArrivedT{Floor: lift.curState.curFloor, ElevatorInd: lift.elInd}
				lift.floorReached <- arrived
			}
		}
	}
}
