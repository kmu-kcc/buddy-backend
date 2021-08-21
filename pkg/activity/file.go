// Package activity provides access to the club activity of the Buddy System.
package activity

import (
	"os"
	"strings"
)

const FilePathPrefix = "./registry"

// File represents a file.
type File string

type Files []File

// Absolute returns the absolute path of f.
func (f File) Absolute() string {
	return strings.Join([]string{FilePathPrefix, string(f)}, "/")
}

// NewFile returns a new file.
func NewFile(filename string) File { return File(strings.TrimSpace(filename)) }

// Delete deletes f.
func (f File) Delete() error {
	return os.Remove(f.Absolute())
}
