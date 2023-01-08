package core

type Direction int
type Floor int
type PassangerStatus int

const (
	Up Direction = iota
	Down
	Stay
)

const (
	ElevatorArrived PassangerStatus = iota
	PassangerArrived
)

type Response struct {
	Direction Direction
	Open      bool
}

type StateBase interface {
	FloorReachedResponse(Floor) (Response, StateBase)
}

type Controller interface {
	NewCall([]StateBase, []Floor, Floor, int) (int, StateBase)
	CabinFloorCall(StateBase, Floor) StateBase
	IntentionToPanelState(Floor, Floor) int
}

type NewPassangerQuery struct {
	CurrentFloor  Floor
	DestinedFloor Floor
	StatusChan    chan PassangerStatus
}
