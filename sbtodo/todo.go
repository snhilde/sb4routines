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
