package sbram

import (
	"errors"
	"strings"
	"io/ioutil"
	"strconv"
	"fmt"
)

// routine is the main object for this package.
// err:     error encountered along the way, if any
// total:   total amount of memory
// total_u: unit of total memory
// used:    amount of memory in current use
// used_u:  unit of used memory
// colors:  trio of user-provided colors for displaying various states
type routine struct {
	err     error
	total   float32
	total_u rune
	used    float32
	used_u  rune
	colors  struct {
		normal  string
		warning string
		error   string
	}
}

// Make and return a new routine object.
func New(colors [3]string) *routine {
	var r routine

	// Do a minor sanity check on the color code.
	for _, color := range colors {
		if !strings.HasPrefix(color, "#") || len(color) != 7 {
			r.err = errors.New("Invalid color")
			return &r
		}
	}
	r.colors.normal  = colors[0]
	r.colors.warning = colors[1]
	r.colors.error   = colors[2]

	return &r
}

// Get the memory resources. Unfortunately, we can't use syscall.Sysinfo() or another syscall function, because it
// doesn't return the necessary information to calculate the actual amount of RAM in use at the moment (namely, it is
// missing the amount of cached RAM). Instead, we're going to read out /proc/meminfo and grab the values we need from
// there. All lines of that file have three fields: field name, value, and unit
func (r *routine) Update() {
	file, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		r.err = err
		return
	}

	total, avail, err := parseFile(string(file))
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

// Format and print the used and total system memory.
func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("%.1f%c/%.1f%c", r.used, r.used_u, r.total, r.total_u)
}

// Parse the meminfo file.
func parseFile(output string) (int, int, error) {
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
	var units = [...]rune{'K', 'M', 'G', 'T', 'P', 'E'}
	var i int

	memory_f := float32(memory)
	for memory_f > 1024 {
		memory_f /= 1024
		i++
	}

	return memory_f, units[i]
}
