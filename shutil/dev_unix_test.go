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

package shutil

import (
	"testing"
)

func TestDev(t *testing.T) {
	devs, err := GetUSBremovables()
	if err != nil {
		t.Fatal(err)
	}

	_, err = FindPartition("foo", devs)
	if err != CmdFindPartError("foo") {
		t.Errorf("FindPartition should get an error")
	}
}
