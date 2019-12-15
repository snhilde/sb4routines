package sbdisk

import (
	"syscall"
	"strings"
	"fmt"
)

// routine is the main object for this package.
// err:   error encountered along the way, if any
// disks: slice of provided filesystems to stat
type routine struct {
	err   error
	disks []fs
}

// fs holds information about a single filesystem.
// path:    given path that will be used to stat the partition
// avail:   used bytes for this filesystem
// avail_u: unit for the used bytes
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

// Copy over the provided filesystem paths and return a new routine object.
func New(paths []string, colors [3]string) *routine {
	var r routine

	for _, path := range paths {
		r.disks = append(r.disks, fs{path: path})
	}

	return &r
}

// For each provided filesystem, get the amounts of used and total disk space and
// convert them into a human-readable size.
func (r *routine) Update() {
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

// Format and print the amounts of disk space for each provided filesystem.
func (r *routine) String() string {
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

// Iteratively decrease the amount of bytes by a step of 2^10 until human-readable.
func shrink(blocks uint64) (uint64, rune) {
	var units = [...]rune{'B', 'K', 'M', 'G', 'T', 'P', 'E'}
	var i int

	for blocks > 1024 {
		blocks >>= 10
		i++
	}

	return blocks, units[i]
}
