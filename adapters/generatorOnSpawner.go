package adapters

import (
	"elevators/generators"
	"time"
)

type CancelHook func()

func (spawner *Spawner) runGenerator(generator generators.Generator, cancel chan bool, period time.Duration) {
	for {
		time.Sleep(period)
		select {
		case <- cancel:
			return
		default:
			from, to := generator.Generate()
			spawner.Spawn(from, to)
		}
	}
}

func (spawner *Spawner) AddAsyncGenerator(generator generators.Generator, period time.Duration) CancelHook {
	cancel := make(chan bool)
	go spawner.runGenerator(generator, cancel, period)
	return func() {
		cancel <- true
	}
}
