package sbvolume

import (
	"os/exec"
	"strings"
	"errors"
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

	return &r
}

func (r *routine) Update() {
	out, err := r.runCmd()
	if err != nil {
		r.err = err
		return
	}

	err = sanityCheck(out)
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

// Make sure the user passed a valid control.
// If the control was not valid, we'll get an error message like this:
// amixer: Unable to find simple control 'bad_control',0
func sanityCheck(out string) error {
	lines := strings.Split(out, "\n")
	if strings.Contains(lines[0], "Unable to find") {
		fields  := strings.Split(lines[0], ":")
		err_msg := strings.TrimSpace(fields[1])
		err_msg  = err_msg[:len(err_msg) - 2]
		return errors.New(err_msg)
	}

	return nil
}
