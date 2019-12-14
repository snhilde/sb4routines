package sbvolume

import (
	"os/exec"
	"strings"
	"errors"
	"strconv"
	"fmt"
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

	// Find the line that has the percentage volume in it.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Playback") && strings.Contains(line, "%]") {
			// We found it. Pull out the volume.
			fields := strings.Fields(line)
			for _, field := range fields {
				if strings.Contains(field, "%]") {
					s        := strings.Trim(field, "[]")
					s         = strings.TrimRight(s, "%")
					vol, err := strconv.Atoi(s)
					if err != nil {
						r.err = err
						return
					}

					r.vol = normalize(vol)
					return
				}
			}
		}
	}

	// If we're here, then we didn't find a volume.
	r.err = errors.New("No volume found for " + r.control)
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("Vol %v%%", r.vol)
}

func (r *routine) runCmd() (string, error) {
	cmd      := exec.Command("amixer", "get", r.control)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// This will ensure that the volumes are multiples of 10 and look nicer.
func normalize(vol int) int {
	return (vol+5) / 10 * 10
}
