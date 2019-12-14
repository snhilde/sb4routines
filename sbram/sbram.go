package sbram

import (
)

type routine struct {
	err error
}

func New() *routine {
	return new(routine)
}

// Unfortunately, we can't use syscall.Sysinfo() or another syscall function, because it doesn't
// return the necessary information to calculate the actual amount of RAM in use at the moment (namely,
// it is missing the amount of cached RAM). Instead, we're going to read out /proc/meminfo and grab
// the values we need from there.
func (r *routine) Update() {
}

func (r *routine) String() string {
	return "ram"
}
