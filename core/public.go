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

type StateBase[T any] interface {
	GetResponse(Floor) Response
	NewFloor(Floor)
	NewCall(Floor, T)
}

type Controller[T any] interface {
	MakeDecision([]StateBase[T], []Floor, Floor, T) (bool, int)
}

type NewPassangerQuery[T any] struct {
	PanelState    T
	CurrentFloor  Floor
	DestinedFloor Floor
	StatusChan    chan PassangerStatus
}
