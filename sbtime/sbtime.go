// Package sbtime displays the current time, according to the format.
package sbtime

import (
	"errors"
	"strings"
	"time"
)

// A routine is the main object for the sbtime package.
// error:  error in colors, if any
// time:   current timestamp
// format: format for displaying time
// colors: trio of user-provided colors for displaying various states
type routine struct {
	err    error
	time   time.Time
	format string
	colors struct {
		normal  string
		warning string
		error   string
	}
}

// Create a new routine object with the current time.
func New(format string, colors [3]string) *routine {
	var r routine

	r.format = format
	r.time   = time.Now()

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

// Update the routine's current time.
func (r *routine) Update() {
	r.time = time.Now()
}

// Print the time in this format: MM D - HH:MM".
func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	if r.time.Second() % 2 == 0 {
		return r.time.Format(r.format)
	} else {
		return r.time.Format(r.format)
	}
}
