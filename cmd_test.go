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
	"errors"
	"testing"
)

var testsOk = []struct {
	cmd string
	out string
	ok  bool
}{
	{"true", "", true},
	{"false", "", false},

	{"ls cmd*.go | wc -l", "2", true},

	{`grep foo shout.go`, "", false},                 // no found
	{`grep package shout.go`, "package shout", true}, // found
}

var testsError = []struct {
	cmd string
	err error // from Stderr
}{
	{"| ls ", errNoCmdInPipe},
	{"| ls | wc", errNoCmdInPipe},
	{"ls|", errNoCmdInPipe},
	{"ls| wc|", errNoCmdInPipe},
	{"ls| |wc", errNoCmdInPipe},

	{"LANG= C find", errEnvVar},
	{"LANG =C find", errEnvVar},

	{`LANG=C find -nop README.md`, errors.New("find: unknown predicate `-nop'")},
}

func TestRun(t *testing.T) {
	//BOOT = true
	//DEBUG = true

	New()
	defer Close()

	for _, v := range testsOk {
		out, ok, _ := Run(v.cmd)

		if out != v.out {
			t.Errorf("`%s` => output got %q, want %q\n", v.cmd, out, v.out)
		}
		if ok != v.ok {
			t.Errorf("`%s` => ok got %t, want %t\n", v.cmd, ok, v.ok)
		}
	}

	for _, v := range testsError {
		_, _, err := Run(v.cmd)
		mainErr := err.(runError).err

		if mainErr.Error() != v.err.Error() {
			t.Errorf("`%s` => error got %q, want %q\n", v.cmd, mainErr, v.err)
		}
	}
}
