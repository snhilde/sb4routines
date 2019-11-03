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
