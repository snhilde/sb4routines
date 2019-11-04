// Package sbtodo prints the first part of a TODO list.
package sbtodo

import (
	"fmt"
	"os"
)

type Routine struct {
}

// Return a new Routine object.
func New() *Routine {
	var r Routine
	return &r
}

 // Read in the given TODO list and format the output according to a few rules:
 //   1. If the file is empty, print "Finished".
 //   2. If the first line has content but the second line is empty, print only the first line.
 //   3. If the first line has content and the second line is indented, print "line1 -> line2".
 //   4. If both lines have content and are both flush, print "line1 | line2".
func (r *Routine) Update() error {
	return nil
}
