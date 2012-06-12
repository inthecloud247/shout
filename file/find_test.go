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
	"os"
	"testing"
)

func TestFind(t *testing.T) {
	defer os.Remove(TEMP_FILE)

	ok, err := ContainString("night", TEMP_FILE)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("ContainString got %t, want %t", ok, !ok)
	}

	ok, err = Contain([]byte("night"), TEMP_FILE)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("Contain got %t, want %t", ok, !ok)
	}
}
