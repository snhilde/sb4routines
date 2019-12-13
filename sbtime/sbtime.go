// Package sbtime displays the current time, according to the format.
package sbtime

import (
	"time"
)

// A routine is the main object for the sbtime package.
// time: current timestamp
type routine struct {
	time time.Time
}

// Create a new routine object and get the current time.
func New() *routine {
	r := routine{time: time.Now()}
	return &r
}

// Update the routine's current time.
func (r *routine) Update() {
	r.time = time.Now()
}

// Print the time in this format: MM D - HH:MM".
func (r *routine) String() string {
	if r.time.Second() % 2 == 0 {
		return r.time.Format("Jan 2 - 03:04")
	} else {
		return r.time.Format("Jan 2 - 03 04")
	}
}
