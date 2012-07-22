// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package system

import (
	"strings"
	"testing"
)

func TestSystem(t *testing.T) {
	sys := NewSystem(Debian)
	cmd := "foo"

	out := sys.Install(cmd)
	if !strings.HasSuffix(out, cmd) {
		t.Error("Install did not get a string")
	}

	if out = sys.Purge(cmd); !strings.HasSuffix(out, cmd) {
		t.Error("Purge did not get a string")
	}
}
