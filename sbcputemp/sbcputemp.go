package sbcputemp

import (
	"errors"
	"strings"
	"fmt"
	"os"
	"io/ioutil"
	"strconv"
)

const (
	COLOR_END = "^d^"
)

// We need to root around in this directory for the device directory for the fan.
const base_dir = "/sys/class/hwmon/"

// routine is the main object for this package.
// err:    error encountered along the way, if any
// path:   path to the device directory, as discovered in findDir()
// files:  slice of files that contain temperature readings
// temp:   average temperature across all sensors, in degrees Celsius
// colors: trio of user-provided colors for displaying various states
type routine struct {
	err      error
	path     string
	files  []os.FileInfo
	temp     int
	colors   struct {
		normal  string
		warning string
		error   string
	}
}

// Find our device directory, build a list of all the temperature sensors in it, and return a new object.
func New(colors [3]string) *routine {
	var r routine

	// Do a minor sanity check on the color code.
	for _, color := range colors {
		if !strings.HasPrefix(color, "#") || len(color) != 7 {
			r.err = errors.New("Invalid color")
			return &r
		}
	}
	r.colors.normal  = "^c" + colors[0] + "^"
	r.colors.warning = "^c" + colors[1] + "^"
	r.colors.error   = "^c" + colors[2] + "^"

	// Error will be handled in Update() and String().
	r.path, r.err = findDir()
	if r.err != nil {
		return &r
	}

	// Error will be handled in Update() and String().
	r.files, r.err = findFiles(r.path)

	return &r
}

// Read out the value of each sensor, get an average of all temperatures, and convert it from milliCelsius to Celsius.
// If we have trouble reading a particular sensor, then we'll skip it on this pass.
func (r *routine) Update() {
	var n int

	if r.path == "" || len(r.files) == 0 {
		return
	}

	r.temp = 0
	for _, file := range r.files {
		b, err := ioutil.ReadFile(r.path + file.Name())
		if err != nil {
			continue
		}

		n, err = strconv.Atoi(strings.TrimSpace(string(b)))
		if err != nil {
			continue
		}

		r.temp += n
	}

	// Get the average temp across all readings.
	r.temp /= len(r.files)

	// Convert to degrees Celsius.
	r.temp /= 1000
}

// Print formatted temperature average in degrees Celsius.
func (r *routine) String() string {
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
		return nil, errors.New("No temperature files")
	}

	return b, nil
}
