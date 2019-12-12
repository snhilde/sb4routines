package sbfan

import (
	"os"
	"io/ioutil"
	"strings"
	"errors"
)

const base_dir = "/sys/class/hwmon"

type Routine struct {
	err  error
	path string
}

func New() *Routine {
	var r Routine

	r.findFile()

	return &r
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return r.path
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
		files, r.err = ioutil.ReadDir(base_dir + "/" + dir.Name())
		if r.err != nil {
			return
		}
		// Find the first file that has a name match. The file we want will start with "fan" and end with "input".
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "fan") && strings.HasSuffix(file.Name(), "input") {
				// We found it.
				r.path = base_dir + dir.Name()
				break;
			}
		}
		if r.path == "" {
			// We found our path. We can stop looking.
			break;
		}
	}

	// Make sure we found something.
	if r.path == "" {
		r.err = errors.New("No fan file")
	}
}
