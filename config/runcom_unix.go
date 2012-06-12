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

// +build !windows

// Package config implements reader of runcom configuration files; with pairs
// key=value.
//
package config

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

var bEq = []byte{'='}

// Loads returns the settings of the configuration file named into a map.
func Load(name string) (map[string]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(file)
	cfg := make(map[string]string)

	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		// Skip comments and blank lines.
		if line[0] == '#' || line[0] == '\n' {
			continue
		}

		fields := bytes.SplitN(line, bEq, 2)
		cfg[string(fields[0])] = string(bytes.Trim(fields[1], `"`)) // remove quotes
	}
	return cfg, file.Close()
}
