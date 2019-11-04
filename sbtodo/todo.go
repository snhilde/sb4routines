// Package sbtodo prints the first part of a TODO list.
package sbtodo

import (
	"fmt"
	"os"
)

// Routine is the main object for this package.
// file will be the open file descriptor for the TODO file.
type Routine struct {
	file *os.File
}

// Return a new Routine object.
// s is the absolute path to the TODO file.
func New(path string) *Routine {
	var r Routine
	var err error

	r.file, err = os.Open(path)
	if err != nil {
		// TODO: handle error
	}

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
