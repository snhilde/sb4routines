package sbram

import (
	"os/exec"
	"strings"
	"errors"
)

type routine struct {
	err error
	total int
	avail int
	used  int
}

func New() *routine {
	return new(routine)
}

// Unfortunately, we can't use syscall.Sysinfo() or another syscall function, because it doesn't
// return the necessary information to calculate the actual amount of RAM in use at the moment (namely,
// it is missing the amount of cached RAM). Instead, we're going to read out /proc/meminfo and grab
// the values we need from there.
func (r *routine) Update() {
	var out []byte

	proc       := exec.Command("cat", "/proc/meminfo")
	out, r.err  = proc.Output()
	if r.err != nil {
		return
	}

	lines := strings.Split(string(out), "\n");
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				r.err = errors.New("Invalid MemTotal fields")
				return
			}

		} else if strings.HasPrefix(line, "MemAvailable") {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				r.err = errors.New("Invalid MemAvailable fields")
				return
			}
		}
	}

	if r.total == 0 || r.avail == 0 {
		r.err = errors.New("Failed to parse out memory fields")
		return
	}
}

func (r *routine) String() string {
	return "ram"
}
