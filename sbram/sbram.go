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
// used:  amount of memory in current use
type routine struct {
	err     error
	total   float32
	total_u rune
	used    float32
	used_u  rune
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
	proc     := exec.Command("cat", "/proc/meminfo")
	out, err := proc.Output()
	if err != nil {
		r.err = err
		return
	}

	total, avail, err := parseCmd(string(out))
	if err != nil {
		r.err = err
		return
	}

	if total == 0 || avail == 0 {
		r.err = errors.New("Failed to parse out memory fields")
		return
	}

	r.total, r.total_u = shrink(total)
	r.used,  r.used_u  = shrink(total - avail)
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "ram"
}

func parseCmd(output string) (int, int, error) {
	var total int
	var avail int
	var err   error

	lines := strings.Split(string(output), "\n");
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				return 0, 0, errors.New("Invalid MemTotal fields")
			}
			total, err = strconv.Atoi(fields[1])
			if err != nil {
				return 0, 0, err
			}

		} else if strings.HasPrefix(line, "MemAvailable") {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				return 0, 0, errors.New("Invalid MemAvailable fields")
			}
			avail, err = strconv.Atoi(fields[1])
			if err != nil {
				return 0, 0, err
			}
		}
	}

	return total, avail, nil
}

// Iteratively decrease the amount of bytes by a step of 2^10 until human-readable.
func shrink(memory int) (float32, rune) {
	var units = [...]rune{'B', 'K', 'M', 'G', 'T', 'P', 'E'}
	var i int

	memory_f := float32(memory)
	for memory_f > 1024 {
		memory_f /= 1024
		i++
	}

	return memory_f, units[i]
}
