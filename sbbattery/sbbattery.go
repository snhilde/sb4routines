package sbbattery

import (
	"os"
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
	r.max = r.readFile("/sys/class/power_supply/BAT0/charge_full")

	return &r
}

// Read current capacity left and calculate percentage from that.
func (r *routine) Update() {
	var now int

	// Handle error reading max capacity.
	if r.max == -1 {
		return
	}

	now = r.readFile("/sys/class/power_supply/BAT0/charge_now")
	if r.err != nil {
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
func (r *routine) readFile(file string) int {
	var f *os.File
	var n int

	f, r.err = os.Open(file)
	if r.err != nil {
		return -1
	}
	defer f.Close()

	_, r.err = fmt.Fscan(f, &n)
	if r.err != nil {
		return -1
	}

	return n
}
