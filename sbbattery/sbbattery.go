package sbbattery

import (
	"io/ioutil"
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
	var val []byte

	val, r.err = ioutil.ReadFile(file)
	if r.err != nil {
		return -1
	}
}
