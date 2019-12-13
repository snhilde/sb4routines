package sbdisk

import (
)

type Routine struct {
	total int
	used  int
	free  int
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "disk"
}
