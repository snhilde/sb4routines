package sbnordvpn

import (
	"os/exec"
	"strings"
	"errors"
)

// routine is the main object in the package.
// err:    error encountered along the way, if any
// b:      buffer to hold connnection string
// colors: trio of user-provided colors for displaying various states
type routine struct {
	err    error
	b      strings.Builder
	colors struct {
		normal  string
		warning string
		error   string
	}
}

// Return a new routine object.
func New(colors [3]string) *routine {
	var r routine

	// Do a minor sanity check on the color code.
	for _, color := range colors {
		if !strings.HasPrefix(color, "#") || len(color) != 7 {
			r.err = errors.New("Invalid color")
			return &r
		}
	}
	r.colors.normal  = colors[0]
	r.colors.warning = colors[1]
	r.colors.error   = colors[2]

	return &r
}

// Run the command and capture the output.
func (r *routine) Update() {
	cmd      := exec.Command("nordvpn", "status")
	out, err := cmd.Output()
	if err != nil {
		r.err = err
		return
	}

	r.parseCommand(string(out))
}

// Format and print the current connection status.
func (r *routine) String() string {
	if r.err != nil {
		return "NordVPN: " + r.err.Error()
	}

	return r.b.String()
}

// Parse the command's output.
func (r *routine) parseCommand(s string) {
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
				fields := strings.Fields(line)
				if len(fields) != 2 {
					r.err = errors.New("Error parsing City")
				} else {
					r.b.Reset()
					r.b.WriteString("Connected: ")
					r.b.WriteString(fields[1])
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
