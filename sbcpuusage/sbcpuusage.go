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
	var r Routine
	return &r
}
