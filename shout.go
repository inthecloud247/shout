// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
