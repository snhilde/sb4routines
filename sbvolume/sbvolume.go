package sbvolume

import (
	"os/exec"
	"strings"
	"errors"
	"strconv"
	"fmt"
)

// routine is the main object for this package.
// err:     error encountered along the way, if any
// control: control to query, as passed in by called
// vol:     system volume, in multiple of ten, as percentage of max
// muted:   true if volume is muted
type routine struct {
	err     error
	control string
	vol     int
	muted   bool
}

// Store the passed-in control value and return a new routine object.
func New(control string, colors [3]string) *routine {
	var r routine

	r.control = control

	return &r
}

// Run the 'amixer' command and parse the output for mute status and volume percentage.
func (r *routine) Update() {
	r.muted = false
	r.vol   = -1

	out, err := r.runCmd()
	if err != nil {
		r.err = err
		return
	}

	// Find the line that has the percentage volume and mute status in it.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Playback") && strings.Contains(line, "%]") {
			// We found it. Check the mute status, then pull out the volume.
			fields := strings.Fields(line)
			for _, field := range fields {
				field = strings.Trim(field, "[]")
				if field == "off" {
					r.muted = true
				} else if strings.HasSuffix(field, "%") {
					s        := strings.TrimRight(field, "%")
					vol, err := strconv.Atoi(s)
					if err != nil {
						r.err = err
						return
					}
					r.vol = normalize(vol)
				}
			}
			break;
		}
	}

	if r.vol < 0 {
		r.err = errors.New("No volume found for " + r.control)
	}
}

// Print either an error, the mute status, or the volume percentage.
func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	if r.muted {
		return "Vol mute"
	}

	return fmt.Sprintf("Vol %v%%", r.vol)
}

// Run the actual 'amixer' command, with the given control.
func (r *routine) runCmd() (string, error) {
	cmd      := exec.Command("amixer", "get", r.control)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Ensure that the volume is a multiple of 10 (so it looks nicer).
func normalize(vol int) int {
	return (vol+5) / 10 * 10
}
