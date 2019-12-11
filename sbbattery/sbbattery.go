package sbbattery

import (
	"io/ioutil"
)

type Routine struct {
	err    error
	charge int
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() error {
	return nil
}

func (r *Routine) String() string {
	return "battery"
}

func (r *Routine) readFile(file string) int {
	var val []byte

	val, r.err = ioutil.ReadFile(file)
	if r.err != nil {
		return -1
	}
}
