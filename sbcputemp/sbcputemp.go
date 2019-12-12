package sbcputemp

import (
)

type Routine struct {
	temp int
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "cpu temp"
}
