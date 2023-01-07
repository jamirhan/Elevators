package core

import (
	"time"
)

func tick(ch chan bool) {
	for {
		time.Sleep(msPerTick * time.Millisecond)
		ch <- true
	}
}


type elevatorContent[T any] struct {
	elevatorState state[T]
	passangers map[Floor][]chan PassangerStatus
}

type Hall[T any] struct {
	states          []elevatorContent[T]
	changeChannels  []chan bool
	newPassanger    chan NewPassangerQuery[T]
	controller      Controller[T]
	floorsMap       map[Floor][]NewPassangerQuery[T]
	elevatorOpened  chan ElevatorArrivedT
}

type ElevatorArrivedT struct {
	Floor       Floor
	ElevatorInd int
}

func CreateHall[T any](elevatorsNum int, initState StateBase[T], newPassanger chan NewPassangerQuery[T], controller Controller[T]) *Hall[T] {
	states := make([]elevatorContent[T], elevatorsNum)
	changeChannels := make([]chan bool, elevatorsNum)
	floorsMap := make(map[Floor][]NewPassangerQuery[T])
	elevatorOpened := make(chan ElevatorArrivedT, 3)
	for i := 0; i < elevatorsNum; i++ {
		states[i].elevatorState.initialize(initState)
		states[i].passangers = make(map[Floor][]chan PassangerStatus)
	}
	for i := 0; i < elevatorsNum; i++ {
		changeChannels[i] = make(chan bool)
		ch := make(chan bool)
		startElevator(changeChannels[i], elevatorOpened, ch, &states[i].elevatorState, i)
		go tick(ch)
	}
	hall := Hall[T]{
		states:          states,
		changeChannels:  changeChannels,
		newPassanger:    newPassanger,
		controller:      controller,
		floorsMap:       floorsMap,
		elevatorOpened:  elevatorOpened,
	}
	return &hall
}

func (hall *Hall[T]) makeDecision(PanelState T, CurrentFloor Floor) (bool, int) {
	stateBases := make([]StateBase[T], len(hall.states))
	curFloors := make([]Floor, len(hall.states))
	for i := 0; i < len(hall.states); i++ {
		hall.states[i].elevatorState.stateMut.Lock()
		stateBases[i] = hall.states[i].elevatorState.state
		curFloors[i] = hall.states[i].elevatorState.curFloor
		hall.states[i].elevatorState.stateMut.Unlock()
	}
	return hall.controller.MakeDecision(stateBases, curFloors, CurrentFloor, PanelState)
}

func (hall *Hall[T]) Routine() {
	for {
		select {
		case newPassanger := <-hall.newPassanger:
			hall.floorsMap[newPassanger.CurrentFloor] = append(hall.floorsMap[newPassanger.CurrentFloor], newPassanger)
			if change_state, index := hall.makeDecision(newPassanger.PanelState, newPassanger.CurrentFloor); change_state {
				hall.states[index].elevatorState.stateMut.Lock()
				hall.states[index].elevatorState.state.NewCall(newPassanger.CurrentFloor, newPassanger.PanelState)
				hall.states[index].elevatorState.stateMut.Unlock()
				hall.changeChannels[index] <- true
			}
		case resp := <-hall.elevatorOpened:
			for _, passanger := range hall.floorsMap[resp.Floor] {
				passanger.StatusChan <- ElevatorArrived
				hall.states[resp.ElevatorInd].elevatorState.stateMut.Lock()
				hall.states[resp.ElevatorInd].elevatorState.state.NewFloor(passanger.DestinedFloor)
				hall.states[resp.ElevatorInd].passangers[passanger.DestinedFloor] = append(hall.states[resp.ElevatorInd].passangers[passanger.DestinedFloor], passanger.StatusChan)
				hall.states[resp.ElevatorInd].elevatorState.stateMut.Unlock()
			}
			for _, passangerChan := range hall.states[resp.ElevatorInd].passangers[resp.Floor] {
				passangerChan <- PassangerArrived
			}
			hall.states[resp.ElevatorInd].passangers[resp.Floor] = make([]chan PassangerStatus, 0)
			hall.floorsMap[resp.Floor] = make([]NewPassangerQuery[T], 0)
		}
	}
}
