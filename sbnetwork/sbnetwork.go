package sbnetwork

import (
	"net"
	"strings"
)

type routine struct {
	err   error
	ilist []net.Interface
}

func New() *routine {
	var r routine

	// Error will be handled in Update() and String().
	r.ilist, r.err = getInterface()
	return &r
}

func (r *routine) Update() {
	if r.err != nil || r.ilist == nil {
		return
	}
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	} else if r.ilist == nil {
		return "No network interfaces found"
	}

	return "network"
}

func getInterface() ([]net.Interface, error) {
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
