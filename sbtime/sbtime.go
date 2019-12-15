// Package sbtime displays the current time, according to the format.
package sbtime

import (
	"time"
)

// A routine is the main object for the sbtime package.
// time:   current timestamp
// format: format for displaying time
type routine struct {
	time   time.Time
	format string
}

// Create a new routine object and get the current time.
func New(format string) *routine {
	var r routine

	r.time   = time.Now()
	r.format = format

	return &r
}

// Update the routine's current time.
func (r *routine) Update() {
	r.time = time.Now()
}

// Print the time in this format: MM D - HH:MM".
func (r *routine) String() string {
	if r.time.Second() % 2 == 0 {
		return r.time.Format(r.format)
	} else {
		return r.time.Format(r.format)
	}
}
