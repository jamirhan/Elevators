package core

import (
	"time"
)

func tick(ch chan bool) {
	for {
		time.Sleep(MsPerTick * time.Millisecond)
		ch <- true
	}
}


type elevatorContent struct {
	elevatorState state
	passangers map[Floor][]chan PassangerStatus
}

type Hall struct {
	states          []elevatorContent
	changeChannels  []chan bool
	newPassanger    chan NewPassangerQuery
	controller      Controller
	floorsMap       map[Floor][]NewPassangerQuery
	elevatorOpened  chan ElevatorArrivedT
}

type ElevatorArrivedT struct {
	Floor       Floor
	ElevatorInd int
}

func CreateHall(elevatorsNum int, initState StateBase, newPassanger chan NewPassangerQuery, controller Controller) *Hall {
	states := make([]elevatorContent, elevatorsNum)
	changeChannels := make([]chan bool, elevatorsNum)
	floorsMap := make(map[Floor][]NewPassangerQuery)
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
	hall := Hall{
		states:          states,
		changeChannels:  changeChannels,
		newPassanger:    newPassanger,
		controller:      controller,
		floorsMap:       floorsMap,
		elevatorOpened:  elevatorOpened,
	}
	return &hall
}

func (hall *Hall) makeDecision(PanelState int, CurrentFloor Floor) (int, StateBase) {
	stateBases := make([]StateBase, len(hall.states))
	curFloors := make([]Floor, len(hall.states))
	for i := 0; i < len(hall.states); i++ {
		hall.states[i].elevatorState.stateMut.Lock()
		stateBases[i] = hall.states[i].elevatorState.state
		curFloors[i] = hall.states[i].elevatorState.curFloor
		hall.states[i].elevatorState.stateMut.Unlock()
	}
	return hall.controller.NewCall(stateBases, curFloors, CurrentFloor, PanelState)
}

func (hall *Hall) Routine() {
	for {
		select {
		case newPassanger := <-hall.newPassanger:
			hall.floorsMap[newPassanger.CurrentFloor] = append(hall.floorsMap[newPassanger.CurrentFloor], newPassanger)
			ind, state := hall.makeDecision(hall.controller.IntentionToPanelState(newPassanger.CurrentFloor, newPassanger.DestinedFloor), newPassanger.CurrentFloor)
			hall.states[ind].elevatorState.stateMut.Lock()
			hall.states[ind].elevatorState.state = state
			hall.states[ind].elevatorState.stateMut.Unlock()
			hall.changeChannels[ind] <- true
		case resp := <-hall.elevatorOpened:
			for _, passanger := range hall.floorsMap[resp.Floor] {
				passanger.StatusChan <- ElevatorArrived
				hall.states[resp.ElevatorInd].elevatorState.stateMut.Lock()
				new_state := hall.controller.CabinFloorCall(hall.states[resp.ElevatorInd].elevatorState.state, passanger.DestinedFloor)
				hall.states[resp.ElevatorInd].elevatorState.state = new_state
				hall.states[resp.ElevatorInd].passangers[passanger.DestinedFloor] = append(hall.states[resp.ElevatorInd].passangers[passanger.DestinedFloor], passanger.StatusChan)
				hall.states[resp.ElevatorInd].elevatorState.stateMut.Unlock()
				hall.changeChannels[resp.ElevatorInd] <- true
			}
			for _, passangerChan := range hall.states[resp.ElevatorInd].passangers[resp.Floor] {
				passangerChan <- PassangerArrived
			}
			hall.states[resp.ElevatorInd].passangers[resp.Floor] = make([]chan PassangerStatus, 0)
			hall.floorsMap[resp.Floor] = make([]NewPassangerQuery, 0)
		}
	}
}
