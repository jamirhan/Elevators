package generators

import "math/rand"
import "elevators/core"

type Random struct {
	MaxFloors int
	Seed int
}

func (generator *Random) Generate() (core.Floor, core.Floor) {
	if generator.MaxFloors <= 1 {
		panic("trying to call generators.Random::Generate() with MaxFloors <= 1")
	}
	rand.NewSource(int64(generator.Seed))
	min := 1
    max := generator.MaxFloors
    first := rand.Intn(max - min) + min
	second := first
	for ;second == first; {
		second = rand.Intn(max - min) + min
	}
	return core.Floor(first), core.Floor(second)
}
