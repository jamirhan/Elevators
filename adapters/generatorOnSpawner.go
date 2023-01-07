package adapters

import (
	"elevators/adapters/generators"
	"time"
)

type CancelHook func()

func (spawner *Spawner[T]) runGenerator(generator generators.Generator[T], cancel chan bool, period time.Duration) {
	for {
		time.Sleep(period)
		select {
		case <- cancel:
			return
		default:
			from, to, panel := generator.Generate()
			spawner.Spawn(from, to, panel)
		}
	}
}

func (spawner *Spawner[T]) AddAsyncGenerator(generator generators.Generator[T], period time.Duration) CancelHook {
	cancel := make(chan bool)
	go spawner.runGenerator(generator, cancel, period)
	return func() {
		cancel <- true
	}
}