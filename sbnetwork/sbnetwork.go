package sbnetwork

import (
	"net"
	"errors"
	"strings"
)

type routine struct {
	err   error
	ilist []net.Interface
}

func New(inames ...string) *routine {
	var r routine

	if len(inames) == 0 {
		// Nothing was passed in. We'll grab the default interfaces.
		// Error will be handled in Update() and String().
		r.ilist, r.err = getInterfaces()
	} else {
		for _, iname := range inames {
			iface, err := net.InterfaceByName(iname)
			if err != nil {
				// Error will be handled in Update() and String().
				r.err = errors.New(iname + ": " + err.Error())
				break;
			}
			r.ilist = append(r.ilist, *iface)
		}
	}

	return &r
}

func (r *routine) Update() {
	if r.err != nil || len(r.ilist) == 0 {
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
