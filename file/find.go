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
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// FindString returns whether the file filename contains the string s. The
// return value is a boolean.
func FindString(s, filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("grep: %s", err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if strings.Contains(line, s) {
			return true, nil
		}
	}
	return false, nil
}
