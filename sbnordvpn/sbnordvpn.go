// Package sbnordvpn displays the current status of the NordVPN connection, including the city and any connection errors.
package sbnordvpn

import (
	"os/exec"
	"strings"
	"errors"
)

var COLOR_END = "^d^"

// routine is the main object in the package.
// err:    error encountered along the way, if any
// b:      buffer to hold connnection string
// color:  current color of the 3 provided
// colors: trio of user-provided colors for displaying various states
type routine struct {
	err    error
	b      strings.Builder
	blink  bool
	color  string
	colors struct {
		normal  string
		warning string
		error   string
	}
}

// Return a new routine object.
func New(colors ...[3]string) *routine {
	var r routine

	// Do a minor sanity check on the color codes.
	if len(colors) == 1 {
		for _, color := range colors[0] {
			if !strings.HasPrefix(color, "#") || len(color) != 7 {
				r.err = errors.New("Invalid color")
				return &r
			}
		}
		r.colors.normal  = "^c" + colors[0][0] + "^"
		r.colors.warning = "^c" + colors[0][1] + "^"
		r.colors.error   = "^c" + colors[0][2] + "^"
	} else {
		// If a color array wasn't passed in, then we don't want to print this.
		COLOR_END = ""
	}

	return &r
}

// Run the command and capture the output.
func (r *routine) Update() {
	var out []byte
	cmd        := exec.Command("nordvpn", "status")
	out, r.err  = cmd.Output()
	if r.err != nil {
		return
	}

	r.parseCommand(string(out))
}

// Format and print the current connection status.
func (r *routine) String() string {
	if r.err != nil {
		return r.colors.error + "NordVPN: " + r.err.Error() + COLOR_END
	}

	return r.color + r.b.String() + COLOR_END
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
	r.b.Reset()
	lines  := strings.Split(s, "\n");
	fields := strings.Fields(lines[0])
	switch fields[1] {
	case "Connected":
		r.color = r.colors.normal
		for _, line := range lines {
			if strings.HasPrefix(line, "City") {
				city := strings.Split(line, ":");
				if len(city) != 2 {
					r.err = errors.New("Error parsing City")
				} else {
					if r.blink {
						r.blink = false
						r.b.WriteString("Connected: ")
					} else {
						r.blink = true
						r.b.WriteString("Connected  ")
					}
					r.b.WriteString(strings.TrimSpace(city[1]))
				}
				break;
			}
		}
	case "Connecting":
		r.color = r.colors.warning
		r.b.WriteString("Connecting...")
	case "Disconnected":
		r.color = r.colors.warning
		r.b.WriteString("Disconnected")
	case "Please check your internet connection and try again.":
		r.err = errors.New("Internet Down")
	default:
		r.err = errors.New(lines[0])
	}
}
