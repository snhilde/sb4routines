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
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() error {
}

func (r *Routine) String() string {
}

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
