package sbcpuusage

import (
	"fmt"
	"os"
	"time"
	"strings"
	"bufio"
)

// Routine is the main object for this package.
type Routine struct {
	err       error
	line      string
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

	r.scanFile()

	return nil
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return "stub"
}

// Open /proc/stat and read out the first line (combined CPU stats) of the file.
func (r *Routine) readFile() {
	var file *os.File

	file, r.err = os.Open("/proc/stat")
	if r.err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	r.line, r.err = reader.ReadString('\n')
}

func (r *Routine) scanFile() {
	// The first line of /proc/stat will look like this:
	// "cpu userVal niceVal sysVal idleVal ..."
	var n int

	n, r.err = fmt.Sscanf(r.line, "cpu %v %v %v %v", &(r.old_stats.user), &(r.old_stats.nice), &(r.old_stats.sys), &(r.old_stats.idle))
	if n != 4 || r.err != nil {
		return
	}
}
