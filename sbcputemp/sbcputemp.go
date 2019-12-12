package sbcputemp

import (
	"fmt"
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
	temp  int
}

func New() *Routine {
	var r Routine

	// Error will be handled in Update() and String().
	r.path, r.err = findDir()
	if r.err != nil {
		return &r
	}

	// Error will be handled in Update() and String().
	r.files, r.err = findFiles(r.path)

	return &r
}

func (r *Routine) Update() {
	var f *os.File
	var n int

	if r.err != nil {
		return
	}

	r.temp = 0
	for _, file := range r.files {
		f, r.err = os.Open(r.path + file.Name())
		if r.err != nil {
			r.temp = 0
			return
		}

		_, r.err = fmt.Fscan(f, &n)
		f.Close()
		if r.err != nil {
			r.temp = 0
			return
		}

		r.temp += n
	}

	// Get the average temp across all readings.
	r.temp /= len(r.files)

	// Convert to degrees Celsius.
	r.temp /= 1000
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("%v °C", r.temp)
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
	var b []os.FileInfo

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "temp") && strings.HasSuffix(file.Name(), "input") {
			// We found a temperature reading. Add it to the list.
			b = append(b, file)
		}
	}

	// Make sure we found something.
	if len(b) == 0 {
		return nil, errors.New("No temp files")
	}

	return b, nil
}
