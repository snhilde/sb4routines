package sbnordvpn

import (
	"os/exec"
	"fmt"
	"strings"
	"errors"
)

// Routine is the main object in the package.
// err: any error encountered along the way, if any
// b:   buffer to hold connnection string
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

// Format and print the current connection status.
func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("NordVPN: %s", r.b.String())
}

// Parse the command's output.
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
	var city string

	lines := strings.Split(s, "\n");
	if lines[0] == "Status: Connected" {
		_, r.err = fmt.Sscanf(lines[3], "City: %s", &city)
		if r.err == nil {
			r.b.Reset()
			r.b.WriteString("Connected: ")
			r.b.WriteString(city)
		}

	} else if lines[0] == "Status: Disconnected" {
		r.err = errors.New("Disconnected")

	} else if lines[0] == "Please check your internet connection and try again." {
		r.err = errors.New("Internet Down")

	} else {
		r.err = errors.New(lines[0])
	}
}
