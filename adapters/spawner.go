package adapters

import (
	"elevators/core"
	"time"
)

type Status int

const (
	WaitingElevator Status = iota
	WaitingFloor
	Arrived
)

type Statistics struct {
	From     core.Floor
	To       core.Floor
	Duration time.Duration
	Status   Status
}

func runQuery[T any](from core.Floor, to core.Floor, pChan chan core.NewPassangerQuery[T], panel T, cancelChan chan bool, statsChan chan Statistics, isOver chan int, id int) {
	statusCh := make(chan core.PassangerStatus)
	query := core.NewPassangerQuery[T]{
		PanelState:    panel,
		CurrentFloor:  from,
		DestinedFloor: to,
		StatusChan:    statusCh,
	}
	pChan <- query
	curStatus := WaitingElevator
	start := time.Now()

L:
	for {
		select {
		case <-cancelChan:
			break L
		case stat := <-statusCh:
			if stat == core.ElevatorArrived {
				curStatus = WaitingFloor
			}
			if stat == core.PassangerArrived {
				curStatus = Arrived
				break L
			}
		}
	}
	isOver <- id
	statsChan <- Statistics{
		From:    from,
		To:      to,
		Duration: time.Since(start),
		Status:  curStatus,
	}
}

type Spawner[T any] struct {
	cancelHook chan bool
	rolling    map[int]chan bool
	isOver     chan int
	pChan      chan core.NewPassangerQuery[T]
	statsChan  chan Statistics
	curId      int
}

func CreateSpawner[T any](pChan chan core.NewPassangerQuery[T]) *Spawner[T] {
	spawner := Spawner[T]{
		cancelHook: make(chan bool),
		rolling:    make(map[int]chan bool),
		isOver:     make(chan int),
		pChan:      pChan,
		statsChan:  make(chan Statistics),
	}
	return &spawner
}

func (spawner *Spawner[T]) Stop() {
	spawner.cancelHook <- true
}

func (spawner *Spawner[T]) Stats() chan Statistics {
	return spawner.statsChan
}

func (spawner *Spawner[T]) Run() {
L:
	for {
		select {
		case <-spawner.cancelHook:
			break L
		case id := <-spawner.isOver:
			delete(spawner.rolling, id)
		}
	}
	for key := range spawner.rolling {
		spawner.rolling[key] <- true
	}
	for len(spawner.rolling) != 0 {
		id := <-spawner.isOver
		delete(spawner.rolling, id)
	}
	close(spawner.statsChan)
}

func (spawner *Spawner[T]) Spawn(from core.Floor, to core.Floor, panel T) {
	cancel := make(chan bool)
	spawner.rolling[spawner.curId] = cancel
	spawner.curId++
	go runQuery(from, to, spawner.pChan, panel, cancel, spawner.statsChan, spawner.isOver, spawner.curId)
}
