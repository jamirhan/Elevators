package generators

import "elevators/core"

type Generator[T any] interface {
	Generate() (core.Floor, core.Floor, T)
}
