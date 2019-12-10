package sbcpuusage

import (
	"fmt"
	"os"
	"time"
	"strings"
	"bufio"
	"errors"
)

// Routine is the main object for this package.
type Routine struct {
	err       error
	old_stats stats
	new_stats stats
}

type stats struct {
	user int
	nice int
	sys  int
	idle int
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() error {
	r.readFile()
	if r.err != nil {
		return r.err
	}

	return nil
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "stub"
}

// Open /proc/stat and read out the CPU stats.
func (r *Routine) readFile() {
	// The first line of /proc/stat will look like this:
	// "cpu userVal niceVal sysVal idleVal ..."
	var file *os.File
	var line  string
	var n     int

	file, r.err = os.Open("/proc/stat")
	if r.err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	line, r.err = reader.ReadString('\n')
	if r.err != nil {
		return
	}

	n, r.err = fmt.Sscanf(line, "cpu %v %v %v %v", &(r.new_stats.user), &(r.new_stats.nice), &(r.new_stats.sys), &(r.new_stats.idle))
	if r.err == nil && n != 4 {
		r.err = errors.New("Failed to read all stats")
	}
}
