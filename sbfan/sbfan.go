package sbfan

import (
	"os"
	"io/ioutil"
)

const base_dir = "/sys/class/hwmon"

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
	var dirs  []os.FileInfo
	var files []os.FileInfo

	// Get all the device directories in the main directory.
	dirs, r.err = ioutil.ReadDir(base_dir)
	if r.err != nil {
		return
	}

	// Search in each device directory to find the fan.
	for _, dir := range dirs {
		files, r.err = ioutil.ReadDir(base_dir + dir.Name())
	}
}
