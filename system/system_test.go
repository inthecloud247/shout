// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package system

import (
	"testing"
)

func TestSystem(t *testing.T) {
	sys := NewSystem(Debian)
	cmd := "postgresql"

	err := sys.Install(cmd)
	if err != nil {
		t.Errorf("\n%s", err)
	}

	if err = sys.Purge(cmd, true); err != nil {
		t.Errorf("\n%s", err)
	}
}