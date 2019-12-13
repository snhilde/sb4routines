package sbdisk

import (
	"golang.org/x/sys/unix"
)

type Routine struct {
	paths []string
}

func New(paths []string) *Routine {
	var r Routine

	r.paths = paths

	return &r
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "disk"
}
