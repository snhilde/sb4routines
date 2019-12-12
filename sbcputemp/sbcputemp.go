package sbcputemp

import (
	"os"
	"ioutil"
)

type Routine struct {
	temp int
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() {
}

func (r *Routine) String() string {
	return "cpu temp"
}

// Find the directory that has the temperature readings. It will be the one with the fan speeds,
// somewhere in /sys/class/hwmon.
func (r *Routine) findDir() string {
	var dirs  []os.FileInfo
	var files []os.FileInfo

	// Get all the device directories in the main directory.
	dirs, r.err = ioutil.ReadDir(base_dir)
	if r.err != nil {
		return ""
	}

	// Search in each device directory to find the fan.
	for _, dir := range dirs {
		path := base_dir + dir.Name() + "/device/"
		files, r.err = ioutil.ReadDir(path)
		if r.err != nil {
			return ""
		}

		// If we encounter a file that matches "fan.*output", then we have the right directory.
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "fan") && strings.HasSuffix(file.Name(), "output") {
				// We found our directory. Return the path.
				return path
			}
		}
	}

	// If we made it here, then we didn't find anything.
	r.err = errors.New("No fan file")
	return ""
}
