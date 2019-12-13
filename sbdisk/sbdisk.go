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

// fs
// path:    given path that will be used to stat the partition
// avail:   available bytes for this filesystem
// avail_u: unit for the available bytes
// total:   total bytes for this filesystem
// total_u: unit for the total bytes
type fs struct {
	path    string
	used    uint64
	used_u  rune
	total   uint64
	total_u rune
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

		total := b.Blocks * uint64(b.Bsize)
		used  := total - (b.Bavail * uint64(b.Bsize))

		r.disks[i].used,  r.disks[i].used_u  = shrink(used)
		r.disks[i].total, r.disks[i].total_u = shrink(total)
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
		fmt.Fprintf(&b, "%s: %v%c/%v%c", r.disks[i].path,
				r.disks[i].used,  r.disks[i].used_u,
				r.disks[i].total, r.disks[i].total_u)
	}

	return b.String()
}

func shrink(blocks uint64) (uint64, rune) {
	var units = [...]rune{'B', 'K', 'M', 'G', 'T'}
	var i int

	for blocks > 1024 {
		blocks >>= 10
		i++
	}

	return blocks, units[i]
}
