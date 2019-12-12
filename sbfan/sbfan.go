package sbfan

import (
	"os"
	"io/ioutil"
	"strings"
	"errors"
)

const base_dir = "/sys/class/hwmon/"

type Routine struct {
	err  error
	path string
	max  int
}

func New() *Routine {
	var r Routine

	// Find the max fan speed file and read its value.
	max_file := r.findFile()
	if r.err == nil {
		r.max := r.readSpeed(max_file)
	}

	return &r
}

func (r *Routine) Update() {
	if r.path == "" {
		return
	}
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return r.path
}

// Find the file that we'll monitor for the fan speed.
// It will be in one of the hardware device directories in /sys/class/hwmon.
func (r *Routine) findFile() os.FileInfo {
	var dirs  []os.FileInfo
	var files []os.FileInfo

	// Get all the device directories in the main directory.
	dirs, r.err = ioutil.ReadDir(base_dir)
	if r.err != nil {
		return nil
	}

	// Search in each device directory to find the fan.
	for _, dir := range dirs {
		path := base_dir + dir.Name() + "/device"
		files, r.err = ioutil.ReadDir(path)
		if r.err != nil {
			return nil
		}

		// Find the first file that has a name match. The file we want will start with "fan" and end with "input".
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "fan") && strings.HasSuffix(file.Name(), "max") {
				// We found it.
				r.path = path
				return file
			}
		}
	}

	// If we made it here, then we didn't find anything.
	r.err = errors.New("No fan file")
	return nil
}

func (r *Routine) readSpeed(file os.FileInfo) int {
}
