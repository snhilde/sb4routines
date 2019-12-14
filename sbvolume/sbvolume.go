package sbvolume

import (
)

type routine struct {
	err error
}

func New() *routine {
	return new(routine)
}

func (r *routine) Update() {
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "volume"
}
