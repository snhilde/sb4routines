package sbcputemp

import (
	"os"
	"io/ioutil"
	"strings"
	"errors"
)

// We need to root around in this directory for the device directory for the fan.
const base_dir = "/sys/class/hwmon/"

type Routine struct {
	err   error
	path  string
	files []os.FileInfo
}

func New() *Routine {
	var r Routine

	r.path, r.err = findDir()
	if r.err != nil {
		return &r
	}

	r.files, r.err = findFiles(r.path)

	return &r
}

func (r *Routine) Update() {
	if r.err != nil {
		return
	}
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "cpu temp"
}

// Find the directory that has the temperature readings. It will be the one with the fan speeds,
// somewhere in /sys/class/hwmon.
func findDir() (string, error) {
	// Get all the device directories in the main directory.
	dirs, err := ioutil.ReadDir(base_dir)
	if err != nil {
		return "", err
	}

	// Search in each device directory to find the fan.
	for _, dir := range dirs {
		path := base_dir + dir.Name() + "/device/"
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return "", err
		}

		// If we encounter a file that matches "fan.*output", then we have the right directory.
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "fan") && strings.HasSuffix(file.Name(), "output") {
				// We found our directory. Return the path.
				return path, nil
			}
		}
	}

	// If we made it here, then we didn't find anything.
	return "", errors.New("No fan file")
}

// Go through given path and build list of files that contain a temperature reading.
// These files will begin with "temp" and end with "input".
func findFiles(path string) ([]os.FileInfo, error) {
}
