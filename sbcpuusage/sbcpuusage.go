package sbcpuusage

import (
	"fmt"
	"os"
	"bufio"
)

// Routine is the main object for this package.
type Routine struct {
	err       error
	old_stats stats
	perc      int
}

type stats struct {
	user int
	nice int
	sys  int
	idle int
}

func New() *Routine {
	var r Routine

	r.readFile(&(r.old_stats))

	return &r
}

func (r *Routine) Update() error {
	var new_stats stats

	r.readFile(&new_stats)
	if r.err != nil {
		return r.err
	}

	used   := (new_stats.user-r.old_stats.user) + (new_stats.nice-r.old_stats.nice) + (new_stats.sys-r.old_stats.sys)
	total  := (new_stats.user-r.old_stats.user) + (new_stats.nice-r.old_stats.nice) + (new_stats.sys-r.old_stats.sys) + (new_stats.idle-r.old_stats.idle)

	// Prevent divide-by-zero error
	if total == 0 {
		r.perc = 0
	} else {
		r.perc = (used * 100) / total
		if r.perc < 0 {
			r.perc = 0
		} else if r.perc > 100 {
			r.perc = 100
		}
	}

	r.old_stats.user = new_stats.user
	r.old_stats.nice = new_stats.nice
	r.old_stats.sys  = new_stats.sys
	r.old_stats.idle = new_stats.idle

	return nil
}

func (r *Routine) String() string {
	if r.err != nil {
		return r.err.Error()
	}

	return fmt.Sprintf("%2d%% CPU", r.perc)
}

// Open /proc/stat and read out the CPU stats.
func (r *Routine) readFile(new_stats *stats) {
	// The first line of /proc/stat will look like this:
	// "cpu userVal niceVal sysVal idleVal ..."
	var file *os.File
	var line  string

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

	// Error will be handled in String().
	_, r.err = fmt.Sscanf(line, "cpu %v %v %v %v", &(new_stats.user), &(new_stats.nice), &(new_stats.sys), &(new_stats.idle))
}
