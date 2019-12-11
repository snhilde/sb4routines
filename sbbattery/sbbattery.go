package sbbattery

import (
	"ioutil"
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

func (r *Routine) readFile(file string) {
	val := ioutil.ReadFile(file)
}
