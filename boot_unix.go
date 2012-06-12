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

// +build linux

package shout

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/kless/Go-Inline/inline"
	"github.com/kless/Go-Term/term"
)

const (
	CMD_WRITE = "/bin/plymouth" // to write during graphical boot
	PATH      = "/sbin:/bin:/usr/sbin:/usr/bin"
)

var USE_CMD_WRITE bool

func init() {
	// Check if there is an the external command to write.
	_, err := os.Stat(CMD_WRITE)
	if os.IsNotExist(err) {
		return
	}

	err = exec.Command(CMD_WRITE, "--ping").Run()
	if _, ok := err.(*exec.ExitError); !ok {

//	if _, ok, _ := shout.Run(CMD_WRITE + " --ping"); ok {
		USE_CMD_WRITE = true
	}
}

// ReadPassword reads a password directly from terminal or through a third program.
func ReadPassword(prompt string) (key []byte, err error) {
	if USE_CMD_WRITE {
		key, err = exec.Command(CMD_WRITE, "ask-for-password", "--prompt="+prompt).Output()
	} else {
		t, err := term.New(syscall.Stdin)
		if err != nil {
			panic(err)
		}
		defer t.Restore()

		t.Echo(true)
		key, err = inline.ReadBytes(prompt)
	}

	if err != nil {
		return nil, fmt.Errorf("ReadPassword: %s", err)
	}
	return
}

// Writef prints a message using the program in CMD_WRITE or to Stderr.
func Writef(format string, a ...interface{}) {
	if USE_CMD_WRITE {
		exec.Command(CMD_WRITE, "message", "--text="+fmt.Sprintf(format, a...)).Run()
	} else {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

// Writefln is like Writef, but adds a new line.
func Writefln(format string, a ...interface{}) { Writef(format+"\n", a) }
