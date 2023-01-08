package main

import (
	"elevators/adapters"
	"elevators/core"
	"elevators/generators"
	"elevators/models/controllers"
	"elevators/models/states"
	"fmt"
	"time"
	// "math/rand"
)

func manualGeneratorExample() {
	newP := make(chan core.NewPassangerQuery)
	controller := controllers.Constant[states.SimpleQueue]{
		Index:    0,
		NewFloor: states.PushBack,
	}
	hall := core.CreateHall(1, states.DefaultSimpleQueue(), newP, &controller)
	go hall.Routine()
	spawner := adapters.CreateSpawner(newP)
	go spawner.Run()
	generator := generators.CreateManual()
	hook := spawner.AddAsyncGenerator(&generator, 1*time.Second)
	go generator.Run()

	i := 0
	for stats := range spawner.Stats() {
		fmt.Println(stats)
		i++
		if i == 10 {
			time.Sleep(10 * time.Second)
			spawner.Stop()
		}
	}
	hook()
}

func randomGeneratorExample() {
	newP := make(chan core.NewPassangerQuery)
	controller := controllers.Constant[states.SimpleQueue]{
		Index:    0,
		NewFloor: states.PushBack,
	}
	hall := core.CreateHall(1, states.DefaultSimpleQueue(), newP, &controller)
	go hall.Routine()
	spawner := adapters.CreateSpawner(newP)
	go spawner.Run()
	generator := generators.Random {
		MaxFloors: 10,
		Seed: 1,
	}
	hook := spawner.AddAsyncGenerator(&generator, 1*time.Second)

	sum := float32(0)
	num := 0
	betweenFloors := core.MsPerTick * core.TicksBetweenFloors
	for stats := range spawner.Stats() {
		sum += float32(stats.Duration.Milliseconds()) / float32(betweenFloors)
		num++
		fmt.Println(stats)
		fmt.Println("Current elevator points:", float32(num) / sum)
	}
	hook()
}

func leastControllerExample() {
	newP := make(chan core.NewPassangerQuery)
	controller := controllers.Least{}
	hall := core.CreateHall(4, states.DefaultSimpleQueue(), newP, &controller)
	go hall.Routine()
	spawner := adapters.CreateSpawner(newP)
	go spawner.Run()
	generator := generators.Random {
		MaxFloors: 10,
		Seed: 1,
	}
	hook := spawner.AddAsyncGenerator(&generator, 1*time.Second)

	sum := float32(0)
	num := 0
	betweenFloors := core.MsPerTick * core.TicksBetweenFloors
	for stats := range spawner.Stats() {
		sum += float32(stats.Duration.Milliseconds()) / float32(betweenFloors)
		num++
		fmt.Println(stats)
		fmt.Println("Current elevator points:", float32(num) / sum)
	}
	hook()
}

func randomControllerExample() {
	newP := make(chan core.NewPassangerQuery)
	controller := controllers.Random{
		MaxFloors: 10,
	}
	hall := core.CreateHall(4, states.DefaultSimpleQueue(), newP, &controller)
	go hall.Routine()
	spawner := adapters.CreateSpawner(newP)
	go spawner.Run()
	generator := generators.Random {
		MaxFloors: 10,
		Seed: 1,
	}
	hook := spawner.AddAsyncGenerator(&generator, 1*time.Second)

	sum := float32(0)
	num := 0
	betweenFloors := core.MsPerTick * core.TicksBetweenFloors
	for stats := range spawner.Stats() {
		sum += float32(stats.Duration.Milliseconds()) / float32(betweenFloors)
		num++
		fmt.Println(stats)
		fmt.Println("Current elevator points:", float32(num) / sum)
	}
	hook()
}

func nonStopExample() {
	newP := make(chan core.NewPassangerQuery)
	controller := controllers.Constant[states.NonStop]{
		Index:    0,
		NewFloor: func(state states.NonStop, floor core.Floor) states.NonStop{
			return state
		},
	}
	hall := core.CreateHall(4, states.NonStop{
		MaxFloors: 10,
	}, newP, &controller)
	go hall.Routine()
	spawner := adapters.CreateSpawner(newP)
	go spawner.Run()
	generator := generators.Random {
		MaxFloors: 10,
		Seed: 1,
	}
	hook := spawner.AddAsyncGenerator(&generator, 1*time.Second)

	sum := float32(0)
	num := 0
	betweenFloors := core.MsPerTick * core.TicksBetweenFloors
	for stats := range spawner.Stats() {
		sum += float32(stats.Duration.Milliseconds()) / float32(betweenFloors)
		num++
		fmt.Println(stats)
		fmt.Println("Current elevator points:", float32(num) / sum)
	}
	hook()
}

func main() {
	nonStopExample()
}
