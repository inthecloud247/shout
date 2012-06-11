// Copyright 2012  The "shout" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package shout handles the shell scripting.
//
// The main tool of this package is the function *Run* which lets to run system
// commands under a new process. It handles pipes, environment variables, and does
// pattern expansion just as in the Bash shell.
//
// The editing of files is very important in the shell scripting to working with
// the configuration files. shout has a great number of functions related to it,
// avoiding to have to use an external command to get the same result, and with the
// advantage of that it is created automatically a backup before of editing a file.
//
package shout

import (
	"os"

	"github.com/kless/shout/boot"
)

var (
	_ENV  []string
	_HOME string // to expand symbol "~"
	BOOT  bool   // does the script is being run during boot?
	DEBUG bool
)

func init() {
	_HOME = os.Getenv("HOME")

	if BOOT {
		_ENV = []string{"PATH=" + boot.PATH}
	} else {
		_ENV = os.Environ()
	}
}
