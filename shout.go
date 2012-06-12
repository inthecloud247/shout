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
// The main function of this package is *Run* which lets to run system commands
// under a new process. It handles pipes, environment variables, and does
// pattern expansion just as in the Bash shell.
//
package shout

import (
	"log"
	//"log/syslog"
	"os"
)

var (
	_ENV  []string
	_HOME string // to expand symbol "~"
	BOOT  bool   // does the script is being run during boot?
	DEBUG bool

	logf *os.File
	log_ *log.Logger
)

func init() {
	_HOME = os.Getenv("HOME")

	if BOOT {
		_ENV = []string{"PATH=" + PATH}
	} else {
		_ENV = os.Environ()
	}
}

// New initializes the log file and set the environment variable PATH.
func New() {
	var err error
	//logFilename := "/tmp/boot.log"
	logFilename := "/var/log/boot_.log"

	if path := os.Getenv("PATH"); path == "" {
		if err := os.Setenv("PATH", PATH); err != nil {
			log.Print(err)
		}
	}

	logf, err = os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Fatal(err)
	}

	log_ = log.New(os.Stderr, "", log.Lshortfile)
}

// Close closes the log file.
func Close() error {
	return logf.Close()
}
