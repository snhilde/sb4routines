package sbdisk

import (
	"syscall"
	"strings"
	"fmt"
)

type Routine struct {
	err   error
	disks []fs
}

// Note: Bavail is the amount of blocks that can actually be used, while
// Bfree is the total amount of unused blocks.
type fs struct {
	path  string
	avail uint64 // unix.Statfs_t.Bavail
	total uint64 // unix.Statfs_t.Blocks
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
		r.disks[i].avail = b.Bavail * uint64(b.Bsize)
		r.disks[i].total = b.Blocks * uint64(b.Bsize)
	}
}

func (r *Routine) String() string {
	var b strings.Builder

	if r.err != nil {
		return r.err.Error()
	}

	for i := range r.disks {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s: %v/%v", r.disks[i].path, r.disks[i].avail, r.disks[i].total)
	}

	return b.String()
}
