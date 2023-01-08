package generators

import "elevators/core"

type Generator interface {
	Generate() (core.Floor, core.Floor)
}
