package adapters

import (
	"elevators/core"
	"fmt"
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

func (stats Statistics) String() string {
	var status string
	switch stats.Status {
	case Arrived:
		status = "Arrived"
	case WaitingElevator:
		status = "WaitingElevator"
	case WaitingFloor:
		status = "WaitingFloor"
	default:
		status = fmt.Sprint(stats.Status)
	}
	return fmt.Sprint("Statistics(from:", stats.From, " to:", stats.To, " duration:", stats.Duration, " status:", status,")")
}

func runQuery(from core.Floor, to core.Floor, pChan chan core.NewPassangerQuery, cancelChan chan bool, statsChan chan Statistics, isOver chan int, id int) {
	statusCh := make(chan core.PassangerStatus)
	query := core.NewPassangerQuery {
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

type Spawner struct {
	cancelHook chan bool
	rolling    map[int]chan bool
	isOver     chan int
	pChan      chan core.NewPassangerQuery
	statsChan  chan Statistics
	curId      int
}

func CreateSpawner(pChan chan core.NewPassangerQuery) *Spawner {
	spawner := Spawner{
		cancelHook: make(chan bool),
		rolling:    make(map[int]chan bool),
		isOver:     make(chan int),
		pChan:      pChan,
		statsChan:  make(chan Statistics),
	}
	return &spawner
}

func (spawner *Spawner) Stop() {
	spawner.cancelHook <- true
}

func (spawner *Spawner) Stats() chan Statistics {
	return spawner.statsChan
}

func (spawner *Spawner) Run() {
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
}

func (spawner *Spawner) Spawn(from core.Floor, to core.Floor) {
	cancel := make(chan bool)
	spawner.rolling[spawner.curId] = cancel
	go runQuery(from, to, spawner.pChan, cancel, spawner.statsChan, spawner.isOver, spawner.curId)
	spawner.curId++
}
