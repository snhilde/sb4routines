package sbnetwork

import (
	"net"
	"errors"
	"strings"
)

type routine struct {
	err     error
	ilist []sbiface
}

type sbiface struct {
	iface     net.Interface
	down_path string
	down      int
	up_path   string
	up        int
}

func New(inames ...string) *routine {
	var r       routine
	var iptr   *net.Interface
	var ilist []net.Interface
	var err     error

	if len(inames) == 0 {
		// Nothing was passed in. We'll grab the default interfaces.
		ilist, err = getInterfaces()
	} else {
		for _, iname := range inames {
			iptr, err = net.InterfaceByName(iname)
			if err != nil {
				// Error will be handled in Update() and String().
				err = errors.New(iname + ": " + err.Error())
				break
			}
			ilist = append(ilist, *iptr)
		}
	}

	// Handle any problems that came up, or build up list of interfaces for later use.
	if err != nil {
		r.err = err
	} else if len(ilist) == 0 {
		r.err = errors.New("No interfaces found")
	} else {
		for _, iface := range ilist {
			down_path := "/sys/class/net/" + iface.Name + "/statistics/rx_bytes"
			up_path   := "/sys/class/net/" + iface.Name + "/statistics/tx_bytes"
			r.ilist = append(r.ilist, sbiface{iface: iface, down_path: down_path, up_path: up_path})
		}
	}

	return &r
}

func (r *routine) Update() {
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "network"
}

func getInterfaces() ([]net.Interface, error) {
	var ilist []net.Interface

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
		ilist = append(ilist, iface)
	}

	return ilist, nil
}

func readFile(path string) (int, error) {
}
