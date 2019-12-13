package sbload

import (
	"syscall"
	"fmt"
)

type routine struct {
	err     error
	load_1  float64
	load_5  float64
	load_15 float64
}

func New() *routine {
	return new(routine)
}

func (r *routine) Update() {
	var info syscall.Sysinfo_t

	r.err = syscall.Sysinfo(&info)
	if r.err != nil {
		return
	}

	r.load_1  = float64(info.Loads[0]) / float64(1 << 16)
	r.load_5  = float64(info.Loads[1]) / float64(1 << 16)
	r.load_15 = float64(info.Loads[2]) / float64(1 << 16)
}

func (r *routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("%.2v %.2v %.2v", r.load_1, r.load_5, r.load_15)
}
