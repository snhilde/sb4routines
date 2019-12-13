package sbdisk

import (
	"syscall"
)

type Routine struct {
	err   error
	disks []fs
}

type fs struct {
	path   string
	bsize  int64 // unix.Statfs_t.Bsize
	btotal int64 // unix.Statfs_t.Blocks
	bavail int64 // unix.Statfs_t.Bavail
	// Note: Bavail is the amount of blocks that can actually be used, while
	// Bfree is the total amount of unused blocks.
}

func New(paths []string) *Routine {
	var r Routine

	for _, path := range paths {
		r.disks = append(r.disks, fs{path: path})
	}

	return &r
}

func (r *Routine) Update() {
	var b syscall.Statfs_t

	for i, disk := range r.disks {
		r.err = syscall.Statfs(disk.path, &b)
		if r.err != nil {
			return
		}
	}
}

func (r *Routine) String() string {
	return "disk"
}
