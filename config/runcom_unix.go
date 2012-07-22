// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
