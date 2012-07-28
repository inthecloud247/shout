// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package config

import (
	"testing"
)

func TestRuncom(t *testing.T) {
	cfg, err := Load("/etc/adduser.conf")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg) == 0 {
		t.Error("cfg got length 0")
	}
}
