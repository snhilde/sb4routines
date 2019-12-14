package sbbattery

import (
	"io/ioutil"
	"strconv"
	"strings"
	"fmt"
)

// Main type for package.
// err:  error encountered along the way, if any
// max:  maximum capacity of battery
// perc: percentage of battery capacity left
type routine struct {
	err  error
	max  int
	perc int
}

// Read the maximum capacity of the battery and return struct.
func New() *routine {
	var r routine

	// Error will be handled in both Update() and String().
	r.max, r.err = readFile("/sys/class/power_supply/BAT0/charge_full")

	return &r
}

// Read current capacity left and calculate percentage from that.
func (r *routine) Update() {
	// Handle error reading max capacity.
	if r.max == -1 {
		return
	}

	now, err := readFile("/sys/class/power_supply/BAT0/charge_now")
	if err != nil {
		r.err = err
		return
	}

	r.perc  = (now * 100) / r.max
	if r.perc < 0 {
		r.perc = 0
	} else if r.perc > 100 {
		r.perc = 100
	}
}

// Print formatted percentage of battery left.
func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("%v%% BAT", r.perc)
}

// Read out value from file.
func readFile(path string) (int, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(strings.TrimSpace(string(b)))
}
