// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shout

import (
	"errors"
	"testing"
)

var testsOk = []struct {
	cmd string
	ok  bool
}{
	// expansion of "~"
	{"ls ~/", true},
}

var testsOutput = []struct {
	cmd string
	out string
	ok  bool
}{
	// values in ok
	{"true", "", true},
	{"false", "", false},
	{`grep foo shout.go`, "", false},                 // no found
	{`grep package cmd.go`, "package shout\n", true}, // found

	// pipes
	{"ls cmd*.go | wc -l", "2\n", true},

	// quotes
	{`sh -c 'echo 123'`, "123\n", true},
	{`sh -c "echo 123"`, "123\n", true},
	{`find -name 'cmd*.go'`, "./cmd.go\n./cmd_test.go\n", true},
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
	for _, v := range testsOk {
		out, ok, _ := Run(v.cmd)

		if ok != v.ok {
			t.Errorf("`%s` => ok got %t, want %t\n", v.cmd, ok, v.ok)
		}

		if string(out) == "" {
			t.Errorf("`%s` => output is empty", v.cmd)
		}
	}

	for _, v := range testsOutput {
		out, ok, _ := Run(v.cmd)

		if string(out) != v.out {
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
