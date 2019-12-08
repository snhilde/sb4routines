// Package sbtodo prints the first part of a TODO list.
package sbtodo

import (
	"os"
	"strings"
	"bufio"
)

// Routine is the main object for this package.
// It contains the data obtained from the specified TODO file, including file info and a copy of the first 2 lines.
type Routine struct {
	path   string
	info   os.FileInfo
	line1  string
	line2  string
	err    error
}

// Return a new Routine object.
// path is the absolute path to the TODO file.
func New(path string) *Routine {
	var r Routine

	r.path = path

	r.info, r.err = os.Stat(path)
	if r.err != nil {
		// We'll print the error in String().
		return &r
	}

	r.readFile()
	if r.err != nil {
		// We'll print the error in String().
		return &r
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
	var new_info os.FileInfo

	new_info, r.err = os.Stat(r.path)
	if r.err != nil {
		return r.err
	}

	// If mtime is not newer than what we already have, we can skip reading the file.
	new_mtime := new_info.ModTime().UnixNano()
	old_mtime := r.info.ModTime().UnixNano()
	if (new_mtime > old_mtime) {
		// The file was modified. Let's parse it.
		r.readFile()
		if r.err != nil {
			return r.err
		}
	}

	r.info = new_info
	return nil
}

func (r *Routine) String() string {
	var b strings.Builder

	// Handle any error we might have received in another stage.
	if r.err != nil {
		return r.err.Error()
	}

	r.line1 = strings.TrimSpace(r.line1)
	if len(r.line1) > 0 {
		// We have content in the first line. Start by adding that.
		b.WriteString(r.line1)
		if len(r.line2) > 0 {
			// We have content in the second line as well. First, let's find out which joiner to use.
			if (strings.HasPrefix(r.line2, "\t")) || (strings.HasPrefix(r.line2, " ")) {
				b.WriteString(" -> ")
			} else {
				b.WriteString(" | ")
			}
			// Next, we'll add the second line.
			b.WriteString(strings.TrimSpace(r.line2))
		}
	} else {
		if len(r.line2) > 0 {
			// We only have a second line. Print just that.
			b.WriteString(strings.TrimSpace(r.line2))
		} else {
			// We don't have content in either line.
			b.WriteString("Finished")
		}
	}

	return b.String()
}

func (r *Routine) readFile() {
	file *os.File

	file, r.err = os.Open(r.path)
	if r.err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	r.line1, r.err = reader.ReadString('\n')
	r.line2, r.err = reader.ReadString('\n')
}
