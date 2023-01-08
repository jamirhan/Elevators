package generators

import (
	"elevators/core"
)

type Query struct {
	from core.Floor
	to core.Floor
}

type Manual struct {
	current []Query
	requestValue chan bool
	newValue chan Query
	generatedValue chan Query
}

func CreateManual() Manual {
	return Manual {
		current: make([]Query, 0),
		requestValue: make(chan bool),
		newValue: make(chan Query),
		generatedValue: make(chan Query),
	}
}

func (generator *Manual) Run() {
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

func (generator *Manual) Generate() (core.Floor, core.Floor) {
	generator.requestValue <- true
	resp := <-generator.generatedValue
	return resp.from, resp.to
}

func (generator Manual) Push(from core.Floor, to core.Floor) {
	generator.newValue <- Query{from, to}
}
