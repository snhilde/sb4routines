package sbfan

import (
	"os"
)

type Routine struct {
	err  error
	path string
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "fan"
}

// Find the file that we'll monitor for the fan speed.
// It will be in one of the hardware device directories in /sys/class/hwmon.
func (r *Routine) findFile() {
	var dirs []os.FileInfo

	// Search in each directory to find the one for the fan.
	dirs, r.err = os.ReadDir("/sys/class/hwmon")
}
