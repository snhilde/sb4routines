package sbload

import (
	"fmt"
)

type routine struct {
	err     error
	load_1  float32
	load_5  float32
	load_15 float32
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

	return fmt.Sprintf("%.1v %.1v %.1v", r.load_1, r.load_5, r.load_15)
}
