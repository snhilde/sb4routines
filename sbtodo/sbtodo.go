// Package sbtodo prints the first part of a TODO list.
package sbtodo

import (
	"fmt"
	"os"
	"time"
	"strings"
	"bufio"
)

// Routine is the main object for this package.
// It contains the data obtained from the specified TODO file, including file info and a copy of the first 2 lines.
// file will be the open file descriptor for the TODO file.
type Routine struct {
	path   string
	info   os.FileInfo
	file  *os.File
	line1  string
	line2  string
}

// Return a new Routine object.
// s is the absolute path to the TODO file.
func New(path string) *Routine {
	var r   Routine
	var err error

	r.path = path

	r.info, err = os.Stat(r.path)
	if err != nil {
		// TODO: handle error
	}

	err = r.readFile()
	if err != nil {
		// TODO: handle error
	}

	return &r
}

 // Read in the given TODO list and format the output according to a few rules:
 //   1. If the file is empty, print "Finished".
 //   2. If the first line has content but the second line is empty, print only the first line.
 //   3. If the first line is empty but the second line has content, print only the second line.
 //   4. If the first line has content and the second line is indented, print "line1 -> line2".
 //   5. If both lines have content and both are flush, print "line1 | line2".
func (r *Routine) Update() error {
	info, err := os.Stat(r.path)
	if err != nil {
		// TODO: handle error
	}

	// If mtime is not newer than what we already have, we can skip reading the file.
	new_mtime := info.ModTime().UnixNano()
	old_mtime := r.info.ModTime().UnixNano()
	if (new_mtime > old_mtime) {
		// The file was modified. Let's parse it.
		err = r.readFile()
		if err != nil {
			// TODO: handle error
		}
	}

	r.info = info
	return nil
}

func (r *Routine) String() string {
	var b strings.Builder

	r.line1 = strings.TrimSpace(r.line1)
	if len(r.line1) > 0 {
		// We have content in the first line. Start by adding that.
		fmt.Fprintf(&b, "%s", r.line1)
		if len(r.line2) > 0 {
			// We have content in the second line as well. Let's find out which conjuction we need.
			if (strings.HasPrefix(r.line2, "\t")) || (strings.HasPrefix(r.line2, " ")) {
				fmt.Fprintf(&b, " -> ")
			} else {
				fmt.Fprintf(&b, " | ")
			}
			fmt.Fprintf(&b, "%s", strings.TrimSpace(r.line2))
		}
	} else {
		if len(r.line2) > 0 {
			// We only have a second line. Print only that.
			fmt.Fprintf(&b, "%s", strings.TrimSpace(r.line2))
		} else {
			// We don't have content in either line.
			fmt.Fprintf(&b, "Finished")
		}
	}

	return b.String()
}

func (r *Routine) readFile() error {
	var err error

	r.file, err = os.Open(r.path)
	if err != nil {
		// TODO: handle error
	}
	defer r.file.Close()

	err = r.readLines()
	if err != nil {
		// TODO: handle error
	}

	return err
}

func (r *Routine) readLines() error {
	var reader *bufio.Reader
	var err     error

	reader = bufio.NewReader(r.file)

	r.line1, err = reader.ReadString('\n')
	if err != nil {
		//TODO: handle error
	}

	r.line2, err = reader.ReadString('\n')
	if err != nil {
		//TODO: handle error
	}

	return err
}
