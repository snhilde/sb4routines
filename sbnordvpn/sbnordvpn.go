package sbnordvpn

import (
	"os/exec"
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
		return "NordVPN: " + r.err.Error()
	}

	return r.b.String()
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
	lines := strings.Split(s, "\n");
	if lines[0] == "Status: Connected" {
		for _, line := range lines {
			if strings.HasPrefix(line, "City") {
				fields := strings.Split(line, ":")
				if len(fields) != 2 {
					r.err = errors.New("Error parsing City")
				} else {
					r.b.Reset()
					r.b.WriteString("Connected: ")
					r.b.WriteString(strings.TrimSpace(fields[1]))
				}
				break;
			}
		}

	} else if lines[0] == "Status: Disconnected" {
		r.err = errors.New("Disconnected")

	} else if lines[0] == "Please check your internet connection and try again." {
		r.err = errors.New("Internet Down")

	} else {
		r.err = errors.New(lines[0])
	}
}
