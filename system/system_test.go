// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package system handles the commands to install and remove programs into an
// operating system.
package system

import (
	"testing"
)

func TestSystem(t *testing.T) {
	sys := NewSystem(Debian)

	out := sys.Install("foo")
	if out == "" {
		t.Error("Install did not get a string")
	}

	if out = sys.Purge("foo"); out == "" {
		t.Error("Purge did not get a string")
	}
}
