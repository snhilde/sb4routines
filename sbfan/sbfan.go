package sbfan

import (
	"os"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
)

const base_dir = "/sys/class/hwmon/"

type Routine struct {
	err      error
	path     string
	max_file os.FileInfo
	out_file os.FileInfo
	max      int
}

func New() *Routine {
	var r Routine

	// Find the max fan speed file and read its value.
	r.findFiles()
	if r.err != nil {
		return &r
	}

	r.max = r.readSpeed(r.max_file)

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
func (r *Routine) findFiles() {
	var dirs   []os.FileInfo
	var files  []os.FileInfo

	// Get all the device directories in the main directory.
	dirs, r.err = ioutil.ReadDir(base_dir)
	if r.err != nil {
		return
	}

	// Search in each device directory to find the fan.
	for _, dir := range dirs {
		path := base_dir + dir.Name() + "/device/"
		files, r.err = ioutil.ReadDir(path)
		if r.err != nil {
			return
		}

		// Find the first file that has a name match. The files we want will start with "fan" and end
		// with "max" or "output".
		prefix := "fan"
		for _, file := range files {
			if strings.HasPrefix(file.Name(), prefix) {
				if strings.HasSuffix(file.Name(), "max") || strings.HasSuffix(file.Name(), "output") {
					// We found one of the two.
					if strings.HasSuffix(file.Name(), "max") {
						r.max_file = file
						prefix     = strings.TrimSuffix(file.Name(), "max")
					} else {
						r.out_file = file
						prefix = strings.TrimSuffix(prefix, "output")
					}
				}

				// If we've found both files, we can stop looking.
				if r.max_file != nil && r.out_file != nil {
					r.path = path
					return
				}
			}
		}
	}

	// If we made it here, then we didn't find anything.
	r.err = errors.New("No fan file")
	return
}

func (r *Routine) readSpeed(file os.FileInfo) int {
	var f *os.File
	var n int

	f, r.err = os.Open(r.path + file.Name())
	if r.err != nil {
		return -1
	}
	defer f.Close()

	_, r.err = fmt.Fscan(f, &n)
	if r.err != nil {
		return -1
	}

	return n
}
