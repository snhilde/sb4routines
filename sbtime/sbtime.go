// Package sbtime displays the current time, according to the format.
package sbtime

import (
	"time"
)

// A Routine is the main object for the sbtime package.
// It holds the current time.
type Routine struct {
	time time.Time
}

// Create a new Routine object.
func New() *Routine {
	r := Routine{time: time.Now()}
	return &r
}

// Update the routine's current time.
func (r *Routine) Update() error {
	r.time = time.Now()

	return nil
}

// Print the time in this format: MM D - HH:MM".
func (r *Routine) String() string {
	if r.time.Second() % 2 == 0 {
		return r.time.Format("Jan 2 - 03:04")
	} else {
		return r.time.Format("Jan 2 - 03 04")
	}
}
