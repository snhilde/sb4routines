package sbbattery

import (
	"os"
	"fmt"
)

type Routine struct {
	err  error
	max  int
	perc int
}

func New() *Routine {
	var r Routine

	// Error will be handled in both Update() and String().
	r.max = r.readFile("/sys/class/power_supply/BAT0/charge_full")

	return &r
}

func (r *Routine) Update() error {
	var now int

	// Handle error reading max capacity.
	if r.max == -1 {
		return r.err
	}

	now = r.readFile("/sys/class/power_supply/BAT0/charge_now")
	if r.err != nil {
		return r.err
	}

	r.perc  = (now * 100) / r.max
	if r.perc < 0 {
		r.perc = 0
	} else if r.perc > 100 {
		r.perc = 100
	}

	return nil
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("%v%% BAT", r.perc)
}

func (r *Routine) readFile(file string) int {
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
