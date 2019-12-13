package sbdisk

import (
	"golang.org/x/sys/unix"
)

type Routine struct {
	err   error
	paths []string
}

func New(paths []string) *Routine {
	var r Routine

	r.paths = paths

	return &r
}

func (r *Routine) Update() {
	var b unix.Statfs_t

	for _, path := range paths {
		r.err = unix.Statfs(path, &b)
		if r.err != nil {
			return
		}
	}
}

func (r *Routine) String() string {
	return "disk"
}
