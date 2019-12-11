package sbnordvpn

import (
)

type Routine struct {
	n int
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "nordvpn"
}
