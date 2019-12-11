package sbnordvpn

import (
	"os/exec"
	"fmt"
	"strings"
)

// Routine is the main object in the package.
// err: any error encountered along the way, if any
type Routine struct {
	err error
}

// Return a new Routine object.
func New() *Routine {
	return new(Routine)
}

// Run the command and capture the output.
func (r *Routine) Update() {
	var out []byte

	proc       := exec.Command("nordvpn", "status")
	out, r.err  = proc.Output()
	if r.err != nil {
		return
	}

	r.parseCommand(string(out))
}

// Format and print the current connection status with this format:
//
func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "nordvpn"
}

func (r *Routine) parseCommand(s string) {
	lines := strings.Split(s, "\n");

	// Possible entries for the first line:
	//     Status: Connected
	//     Status: Disconnected
	//     Please check your internet connection and try again.
}
