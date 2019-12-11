package sbbattery

import (
	"os"
	"fmt"
	"errors"
)

type Routine struct {
	err    error
	charge int
	max    int
}

func New() *Routine {
	var r Routine

	// Error will be handled in both Update() and String().
	r.max = r.readFile("/sys/class/power_supply/BAT0/charge_full")

	return &r
}

func (r *Routine) Update() error {
	// Handle error reading max capacity.
	if r.max == -1 {
		return r.err
	}

	return nil
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "battery"
}

func (r *Routine) readFile(file string) int {
	var f *os.File
	var i int
	var n int

	f, r.err = os.Open(file)
	if r.err != nil {
		return -1
	}

	i, r.err = fmt.Fscan(f, &n)
	if i != 1 || r.err != nil {
		if r.err == nil {
			r.err = errors.New("Failed to read file")
		}
		return -1
	}

	return n
}
