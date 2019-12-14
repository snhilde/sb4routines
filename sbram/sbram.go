package sbram

import (
	"os/exec"
	"strings"
	"errors"
	"strconv"
)

// routine is the main object for this package.
// err:   error encountered along the way, if any
// total: total amount of memory
// avail: amount of memory available
// used:  amount of memory in current use
type routine struct {
	err error
	total int
	avail int
	used  int
}

// Make and return a new routine object.
func New() *routine {
	return new(routine)
}

// Get the memory resources. Unfortunately, we can't use syscall.Sysinfo() or another syscall function, because it
// doesn't return the necessary information to calculate the actual amount of RAM in use at the moment (namely, it is
// missing the amount of cached RAM). Instead, we're going to read out /proc/meminfo and grab the values we need from
// there. All lines of that file have three fields: field name, value, and unit
func (r *routine) Update() {
	var out []byte

	proc       := exec.Command("cat", "/proc/meminfo")
	out, r.err  = proc.Output()
	if r.err != nil {
		return
	}

	r.total, r.avail, r.err = parseCmd(string(out))
	if r.err != nil {
		return
	}

	if r.total == 0 || r.avail == 0 {
		r.err = errors.New("Failed to parse out memory fields")
		return
	}

}

func (r *routine) String() string {
	return "ram"
}

func parseCmd(output string) (int, int, error) {
	var total int
	var avail int

	lines := strings.Split(string(out), "\n");
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				return 0, 0, errors.New("Invalid MemTotal fields")
			}
			total = strconv.Atoi(fields[1])

		} else if strings.HasPrefix(line, "MemAvailable") {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				return 0, 0, errors.New("Invalid MemAvailable fields")
			}
			avail = strconv.Atoi(fields[1])
		}
	}

	return total, avail, nil
}
