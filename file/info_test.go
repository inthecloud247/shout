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

package file

import (
	"testing"
)

func TestInfo(t *testing.T) {
	ok, err := IsDir("../file")
	if err != nil {
		t.Error(err)
	} else if !ok {
		t.Error("IsDir got false")
	}

	fi, err := NewInfo("info.go")
	if err != nil {
		t.Fatal(err)
	}

	if !fi.IsModer() {
		t.Error("IsModer got false")
	}
	if !fi.IsModew() {
		t.Error("IsModew got false")
	}
	if fi.IsModex() {
		t.Error("IsModex got true")
	}
	if !fi.IsMode(0600) { // rw
		t.Error("IsMode got false")
	}
}
