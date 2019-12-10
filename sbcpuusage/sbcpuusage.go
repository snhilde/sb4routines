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
	err   error
	stats string
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() error {
}

func (r *Routine) String() string {
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

	r.stats, r.err = reader.ReadString('\n')
}
