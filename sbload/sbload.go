package sbload

import (
	"errors"
	"strings"
	"syscall"
	"fmt"
)

// routine is the main object in the package.
// err:     error encountered along the way, if any
// load_1:  load average over the last second
// load_5:  load average over the last 5 seconds
// load_15: load average over the last 15 seconds
// colors:  trio of user-provided colors for displaying various states
type routine struct {
	err     error
	load_1  float64
	load_5  float64
	load_15 float64
	colors  struct {
		normal  string
		warning string
		error   string
	}
}

// Return a new rountine object.
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

// Call Sysinfo() method and calculate load averages.
func (r *routine) Update() {
	var info syscall.Sysinfo_t

	r.err = syscall.Sysinfo(&info)
	if r.err != nil {
		return
	}

	// Each load average must be divided by 2^16 to get the same format as /proc/loadavg.
	r.load_1  = float64(info.Loads[0]) / float64(1 << 16)
	r.load_5  = float64(info.Loads[1]) / float64(1 << 16)
	r.load_15 = float64(info.Loads[2]) / float64(1 << 16)
}

// Print the 3 load averages with 2 decimal places of precision.
func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("%.2f %.2f %.2f", r.load_1, r.load_5, r.load_15)
}
