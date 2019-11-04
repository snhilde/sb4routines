package sbdelim

import (
	"time"
)

// A Routine is the main object for the sbdelim package.
// For this package, it will always be set to ";".
type Routine string

// Create a new Routine object.
func New() *Routine {
	var r Routine
	return &r
}

// Do nothing, because we never want to change this.
func (r *Routine) Update() error {
	return nil
}

// Print the delimiter for the dualstatus patch.
func (r *Routine) String() string {
	return ";"
}

// Sleep for a long time.
func (r *Routine) Sleep(d time.Duration) {
	time.Sleep(time.Hour)
}
