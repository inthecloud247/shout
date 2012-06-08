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

package shout

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// PATH to use in boot scripts.
const (
	_PATH = "/sbin:/bin:/usr/sbin:/usr/bin"
)

// == Errors
// ==
var (
	errEnvVar      = errors.New("the format of the variable has to be VAR=value")
	errNoCmdInPipe = errors.New("no command around of pipe")
)

type extraCmdError string

func (e extraCmdError) Error() string {
	return "command not added to " + string(e)
}

type runError struct {
	cmd     string
	debug   string
	errType string
	err     error
}

func (e runError) Error() string {
	if e.debug != "" {
		e.debug = "\n\tDEBUG: " + e.debug
	}
	return fmt.Sprintf("[Run] `%s`%s\n\t%s: %s", e.cmd, e.debug, e.errType, e.err)
}

// * * *

// Run executes external commands with access to shell features such as filename
// wildcards, shell pipes and environment variables.
//
// This function avoids to have execute commands through a shell since an
// unsanitized input from an untrusted source makes a program vulnerable to
// shell injection, a serious security flaw which can result in arbitrary
// command execution.
//
// The most of commands return a text in output or an error if any. ok is used
// in commands like *grep*, *find*, or *cmp* to indicate if the serach is matched.
func Run(command string) (output string, ok bool, err error) {
	var (
		env            []string
		cmds           []*exec.Cmd
		outPipes       []io.ReadCloser
		stdout, stderr bytes.Buffer
	)

	if BOOT {
		env = []string{"PATH=" + _PATH}
	} else {
		env = os.Environ()
	}

	commands := strings.Split(command, "|")
	lastIdxCmd := len(commands) - 1

	// Check lonely pipes.
	for _, cmd := range commands {
		if strings.TrimSpace(cmd) == "" {
			err = runError{command, "", "ERR", errNoCmdInPipe}
			return
		}
	}

	for i, cmd := range commands {
		cmdEnv := env  // evironment variables for each command
		indexArgs := 1 // position where the arguments start
		fields := strings.Fields(cmd)
		lastIdxFields := len(fields) - 1

		// == Get environment variables in the first arguments, if any.
		for j, fCmd := range fields {
			if fCmd[len(fCmd)-1] == '=' || // VAR= foo
				(j < lastIdxFields && fields[j+1][0] == '=') { // VAR =foo
				err = runError{command, "", "ERR", errEnvVar}
				return
			}

			if strings.ContainsRune(fields[0], '=') {
				cmdEnv = append([]string{fields[0]}, env...) // Insert the environment variable
				fields = fields[1:]                          // and it is removed from arguments
			} else {
				break
			}
		}
		// ==

		cmdPath, e := exec.LookPath(fields[0])
		if e != nil {
			err = runError{command, "", "ERR", e}
			return
		}

		// == Get the path of the next command, if any
		for j, fCmd := range fields {
			cmdBase := path.Base(fCmd)

			if cmdBase != "sudo" && cmdBase != "xargs" {
				break
			}
			// It should have an extra command.
			if j+1 == len(fields) {
				err = runError{command, "", "ERR", extraCmdError(cmdBase)}
				return
			}

			nextCmdPath, e := exec.LookPath(fields[j+1])
			if e != nil {
				err = runError{command, "", "ERR", e}
				return
			}

			if fields[j+1] != nextCmdPath {
				fields[j+1] = nextCmdPath
				indexArgs = j + 2
			}
		}

		// == Expand the shell file name pattern in arguments, if any
		expand := make(map[int][]string, len(fields))

		for j := indexArgs; j < len(fields); j++ {
			// Skip flags
			if fields[j][0] == '-' {
				continue
			}

			names, e := filepath.Glob(fields[j])
			if e != nil {
				err = runError{command, "", "ERR", e}
				return
			}
			if names != nil {
				expand[j] = names
			}
		}

		// Substitute the names generated for the pattern starting from last field.
		if len(expand) != 0 {
			for j := len(fields) - indexArgs; j >= indexArgs; j-- {
				if v, ok := expand[j]; ok {
					fields = append(fields[:j], append(v, fields[j+1:]...)...)
				}
			}
		}

		// == Create command
		c := &exec.Cmd{
			Path: cmdPath,
			Args: append([]string{fields[0]}, fields[1:]...),
			Env:  cmdEnv,
		}

		// == Connect pipes
		outPipe, e := c.StdoutPipe()
		if e != nil {
			err = runError{command, "", "ERR", e}
			return
		}

		if i == 0 {
			c.Stdin = os.Stdin
		} else {
			c.Stdin = outPipes[i-1] // anterior output
		}

		// == Buffers
		c.Stderr = &stderr

		// Only save the last output
		if i == lastIdxCmd {
			c.Stdout = &stdout
		}

		// == Start command
		if e := c.Start(); e != nil {
			err = runError{command,
				fmt.Sprintf("Path: %s | Args: %s", c.Path, c.Args),
				"Start", fmt.Errorf("%s", c.Stderr)}
			return
		}

		// ==
		cmds = append(cmds, c)
		outPipes = append(outPipes, outPipe)
	}

	for _, c := range cmds {
		if e := c.Wait(); e != nil {
			_, isExitError := e.(*exec.ExitError)

			// Error type due I/O problems.
			if !isExitError {
				err = runError{command,
					fmt.Sprintf("Path: %s | Args: %s", c.Path, c.Args),
					"Wait", fmt.Errorf("%s", c.Stderr)}
				return
			}

			if c.Stderr != nil {
				if stderr := fmt.Sprintf("%s", c.Stderr); stderr != "" {
					stderr = strings.TrimRight(stderr, "\n")
					err = runError{command,
						fmt.Sprintf("Path: %s | Args: %s", c.Path, c.Args),
						"Stderr", fmt.Errorf("%s", stderr)}
					return
				}
			}
		} else {
			ok = true
		}
	}

	return strings.Trim(stdout.String(), "\n"), ok, nil
}

// Sudo calls to command sudo.
// If anything command needs to use sudo, then could be used this function at
// the beginning so there is not to wait until that it been requested later.
func Sudo() {
	out, _, _ := Run("sudo /bin/true")
	println("out:", out)
}
