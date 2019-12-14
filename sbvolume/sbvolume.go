package sbvolume

import (
	"os/exec"
)

type routine struct {
	err     error
	control string
}

func New(control string) *routine {
	var r routine

	r.control = control

	return &r
}

func (r *routine) Update() {
	cmd      := exec.Command("amixer", "get", r.control)
	out, err := cmd.Output()
	if err != nil {
		r.err = err
		return
	}
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "volume"
}
