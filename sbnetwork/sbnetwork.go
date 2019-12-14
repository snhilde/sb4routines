package sbnetwork

import (
	"net"
	"errors"
	"strings"
	"os"
	"fmt"
)

type routine struct {
	err     error
	ilist []sbiface
}

type sbiface struct {
	name      string
	down_path string
	up_path   string
	down_old  int
	up_old    int
	down_new  int
	up_new    int
}

func New(inames ...string) *routine {
	var r       routine
	var ilist []string
	var err     error

	if len(inames) == 0 {
		// Nothing was passed in. We'll grab the default interfaces.
		ilist, err = getInterfaces()
	} else {
		for _, iname := range inames {
			// Make sure we have a valid interface name.
			_, err = net.InterfaceByName(iname)
			if err != nil {
				// Error will be handled in Update() and String().
				err = errors.New(iname + ": " + err.Error())
				break
			}
			ilist = append(ilist, iname)
		}
	}

	// Handle any problems that came up, or build up list of interfaces for later use.
	if err != nil {
		r.err = err
	} else if len(ilist) == 0 {
		r.err = errors.New("No interfaces found")
	} else {
		for _, iname := range ilist {
			down_path := "/sys/class/net/" + iname + "/statistics/rx_bytes"
			up_path   := "/sys/class/net/" + iname + "/statistics/tx_bytes"
			r.ilist = append(r.ilist, sbiface{name: iname, down_path: down_path, up_path: up_path})
		}
	}

	return &r
}

func (r *routine) Update() {
	for i, iface := range r.ilist {
		r.ilist[i].down_old = iface.down_new
		r.ilist[i].up_old   = iface.up_new

		down, err := readFile(iface.down_path)
		if err != nil {
			r.err = err
			break
		}
		r.ilist[i].down_new = down

		up, err := readFile(iface.up_path)
		if err != nil {
			r.err = err
			break
		}
		r.ilist[i].up_new = up
	}
}

func (r *routine) String() string {
	var b strings.Builder

	if r.err != nil {
		return r.err.Error()
	}

	for i, iface := range r.ilist {
		down, down_u := shrink(iface.down_new - iface.down_old)
		up, up_u     := shrink(iface.up_new   - iface.up_old)

		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s: %4v%c↓/%4v%c↑", iface.name, down, down_u, up, up_u)
	}

	return b.String()
}

func getInterfaces() ([]string, error) {
	var inames []string

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Name == "lo" {
			// Skip loopback.
			continue
		} else if !strings.Contains(iface.Flags.String(), "up") {
			// If the network is not up, then we don't need to monitor it.
			continue
		}
		inames = append(inames, iface.Name)
	}

	return inames, nil
}

func readFile(path string) (int, error) {
	var n int

	f, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	_, err = fmt.Fscan(f, &n)
	if err != nil {
		return -1, err
	}

	return n, nil
}

// Iteratively decrease the amount of bytes by a step of 2^10 until human-readable.
func shrink(bytes int) (int, rune) {
	var units = [...]rune{'B', 'K', 'M', 'G', 'T', 'P', 'E'}
	var i int

	for bytes > 1024 {
		bytes >>= 10
		i++
	}

	return bytes, units[i]
}
