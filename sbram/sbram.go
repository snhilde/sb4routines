package sbram

import (
)

type routine struct {
}

func New() *routine {
	return new(routine)
}

func (r *routine) Update() {
}

func (r *routine) String() string {
	return "ram"
}
