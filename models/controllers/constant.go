package controllers

import "elevators/core"

type Constant[T any] struct {
	index int
}

func (controller Constant[T]) MakeDecision([]core.StateBase[T], []core.Floor, core.Floor, T) (bool, int) {
	return true, controller.index
}

func ConstantDecision[T any](index int) Constant[T] {
	return Constant[T]{index: index}
}
