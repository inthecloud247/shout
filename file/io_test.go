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
	"path/filepath"
	"testing"
)

func TestBackupSuffix(t *testing.T) {
	okFilenames := []string{"foo+1~", "foo+2~", "foo+5~", "foo+8~", "foo+9~"}
	badFilenames := []string{"foo+0~", "foo+10~", "foo+11~", "foo+22~"}

	for _, v := range okFilenames {
		ok, err := filepath.Match("foo"+_BACKUP_SUFFIX, v)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Errorf("%q should be matched", v)
		}
	}

	for _, v := range badFilenames {
		ok, err := filepath.Match("foo"+_BACKUP_SUFFIX, v)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Errorf("%q should not be matched", v)
		}
	}
}
