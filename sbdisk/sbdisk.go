package sbdisk

import (
	"syscall"
)

type Routine struct {
	err   error
	disks []fs
}

// Note: Bavail is the amount of blocks that can actually be used, while
// Bfree is the total amount of unused blocks.
type fs struct {
	path  string
	avail int64 // unix.Statfs_t.Bavail
	total int64 // unix.Statfs_t.Blocks
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
		r.disks[i].avail = b.Bavail
		r.disks[i].total = b.Blocks
	}
}

func (r *Routine) String() string {
	return "disk"
}
