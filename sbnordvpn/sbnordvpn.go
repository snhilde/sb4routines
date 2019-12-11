package sbnordvpn

import (
	"os/exec"
	"fmt"
	"strings"
	"errors"
)

// Routine is the main object in the package.
// err: any error encountered along the way, if any
type Routine struct {
	err error
	b   strings.Builder
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
	// If there is a connection to the VPN, the command will return this format:
	//     Status: Connected
	//     Current server: <server.url>
	//     Country: <country>
	//     City: <city>
	//     Your new IP: <the.new.IP.address>
	//     Current technology: <tech>
	//     Current protocol: <protocol>
	//     Transfer: <bytes> <unit> received, <bytes> <unit> sent
	//     Uptime: <human-readable time>
	//
	// If there is no connection, the command will return this:
	//     Status: Disconnected
	//
	// If there is no Internet connection, the command will return this:
	//     Please check your internet connection and try again.
	r.b.Reset()

	lines := strings.Split(s, "\n");
	if lines[0] == "Status: Connected" {
	} else if lines[0] == "Status: Disconnected" {
		r.b.WriteString("Disconnected")
	} else if lines[0] == "Please check your internet connection and try again." {
		r.b.WriteString("Internet Down")
	} else {
		r.err = errors.New(lines[0])
	}
}
