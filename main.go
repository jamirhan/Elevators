package main

import (
	"elevators/adapters"
	"elevators/adapters/generators"
	"elevators/core"
	"elevators/models/controllers"
	"elevators/models/panels"
	"elevators/models/states"
	"fmt"
	"time"
)

func main() {
	newP := make(chan core.NewPassangerQuery[panels.OneButton])
	hall := core.CreateHall[panels.OneButton](1, states.DefaultSimpleQueue[panels.OneButton](), newP, controllers.ConstantDecision[panels.OneButton](0))
	go hall.Routine()
	spawner := adapters.CreateSpawner(newP)
	go spawner.Run()
	generator := generators.CreateManual(panels.ToPanel)
	hook := spawner.AddAsyncGenerator(&generator, 1 * time.Second)

	generator.Push(1, 2)
	generator.Push(2, 4)
	generator.Push(3, 5)

	i := 0
	for stats := range spawner.Stats() {
		fmt.Println(stats)
		i++
		if i >= 10 {
			break
		}
	}
	hook()
	spawner.Stop()
}
