package sbtime

import (
	"time"
)

// A Routine is the main object for the sbtime package.
type Routine struct {
	time time.Time
}

// Create a new Routine object
func New() *Routine {
	r := Routine{time: time.Now()}
	return &r
}

// Update the routine's current time
func (r *Routine) Update() error {
	r.time = time.Now()

	return nil
}
