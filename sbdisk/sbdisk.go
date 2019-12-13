package sbdisk

import (
	"golang.org/x/sys/unix"
)

type Routine struct {
	paths []string
}

func New(paths []string) *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "disk"
}
