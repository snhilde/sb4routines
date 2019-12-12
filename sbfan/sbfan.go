package sbfan

import (
)

type Routine struct {
	path string
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "fan"
}

// Find the file that we'll monitor for the fan speed.
func (r *Routine) findFile() {
}
