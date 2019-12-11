package sbnordvpn

import (
	"os/exec"
	"fmt"
)

type Routine struct {
	err error
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
	var out []byte

	proc       := exec.Command("nordvpn", "status")
	out, r.err  = proc.Output()
	if r.err != nil {
		return
	}
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "nordvpn"
}
