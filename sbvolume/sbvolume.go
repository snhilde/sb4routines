package sbvolume

import (
	"os/exec"
)

type routine struct {
	err     error
	control string
	vol     int
	mute    bool
}

func New(control string) *routine {
	var r routine

	r.control = control

	// Make sure the user passed a valid control.
	out, err := r.runCmd()
	if err != nil {
		r.err = err
		return &r
	}

	// If the control was not valid, we'll get an error message like this:
	// amixer: Unable to find simple control 'bad_control',0

	return &r
}

func (r *routine) Update() {
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

func (r *routine) runCmd() (string, error) {
	cmd      := exec.Command("amixer", "get", r.control)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
