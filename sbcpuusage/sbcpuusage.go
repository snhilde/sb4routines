package sbcpuusage

import (
	"os/exec"
	"strings"
	"errors"
	"strconv"
	"fmt"
	"os"
	"bufio"
)

// Routine is the main object for this package.
// err:       error encountered along the way, if any
// old_stats: CPU stats from last read
// perc:      percentage of CPU currently being used
type Routine struct {
	err       error
	threads   int
	old_stats stats
	perc      int
}

// Type to hold values of different CPU stats
type stats struct {
	user int
	nice int
	sys  int
	idle int
}

// Get current CPU stats and return Routine object.
func New() *Routine {
	var r Routine

	r.threads, r.err = numThreads()
	if r.err != nil {
		return &r
	}

	r.readFile(&(r.old_stats))

	return &r
}

// Get current CPU stats, compare to last-read stats, and calculate percentage of CPU being used.
func (r *Routine) Update() {
	var new_stats stats

	r.readFile(&new_stats)
	if r.err != nil {
		return
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
}

// Print formatted CPU percentage.
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

// The shell command 'lscpu' will return a variety of CPU information, including the number of threads
// per CPU core. We don't care about the number of cores, because we're already reading in the
// averaged total. We only want to know if we need to be changing its range. To get this number, we're
// going to loop through each line of the output until we find "Thread(s) per socket".
func numThreads() (int, error) {
	proc     := exec.Command("lscpu")
	out, err := proc.Output()
	if err != nil {
		return -1, err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Thread(s) per core") {
			fields := strings.Split(line, ":")
			if len(fields) != 2 {
				return -1, errors.New("Invalid fields")
			}
			s      := strings.TrimSpace(fields[1])
			n, err := strconv.Atoi(s)
			if err != nil {
				return -1, err
			}
			return n, nil
		}
	}

	// If we made it this far, then we didn't find anything.
	return -1, errors.New("Failed to find number of threads")
}
