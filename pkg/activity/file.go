// Copyright 2021 KMU KCC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package activity provides access to the club activity of the Buddy System.
package activity

import (
	"os"
	"strings"
)

var (
	home, _        = os.UserHomeDir()
	FilePathPrefix = strings.Join([]string{home, "registry"}, "/")
)

// File represents a file.
type File string

type Files []File

// Absolute returns the absolute path of f.
func (f File) Absolute() string { return strings.Join([]string{FilePathPrefix, string(f)}, "/") }

// NewFile returns a new file.
func NewFile(filename string) File { return File(strings.TrimSpace(filename)) }

// Delete deletes f.
func (f File) Delete() error { return os.Remove(f.Absolute()) }
