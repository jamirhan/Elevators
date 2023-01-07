package generators

import (
	"elevators/core"
	"fmt"
)

type Query struct {
	from core.Floor
	to core.Floor
}

type Manual[T any] struct {
	toPanel func(core.Floor, core.Floor) T
	current []Query
	requestValue chan bool
	newValue chan Query
	generatedValue chan Query
}

func CreateManual[T any](toPanel func(core.Floor, core.Floor) T) Manual[T] {
	return Manual[T] {
		toPanel: toPanel,
		current: make([]Query, 0),
		requestValue: make(chan bool),
		newValue: make(chan Query),
		generatedValue: make(chan Query),
	}
}

func (generator *Manual[T]) Run() {
	waiting := 0
	for {
		select {
		case <-generator.requestValue:
			waiting++
		case value := <- generator.newValue:
			generator.current = append(generator.current, value)
		default:
			if waiting > 0 && len(generator.current) > 0 {
				generator.generatedValue <- generator.current[0]
				generator.current = generator.current[1:]
				waiting--
			}
		}
	}
}

func (generator *Manual[T]) Generate() (core.Floor, core.Floor, T) {
	fmt.Println("generate()")
	generator.requestValue <- true
	resp := <-generator.generatedValue
	fmt.Println("got value")
	return resp.from, resp.to, generator.toPanel(resp.from, resp.to)
}

func (generator Manual[T]) Push(from core.Floor, to core.Floor) {
	generator.newValue <- Query{from, to}
}
