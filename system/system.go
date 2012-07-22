// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package system handles the commands to install and remove programs into an
// operating system.
package system

type pkg interface {
	Install(s string) string
	Remove(s string) string
	Purge(s string) string
}

// system represents the operating system or Linux flavor.
type system int

// Operating system or Linux flavor.
const (
	Debian system = iota + 1
)

// * * *

type debian system

// Install returns the command to install a program.
func (d debian) Install(s string) string {
	return "apt-get install -y " + s
}

// Remove returns the command to remove a program.
func (d debian) Remove(s string) string {
	return "apt-get remove -y " + s
}

// Purge returns the command to purge a program.
func (d debian) Purge(s string) string {
	return "apt-get purge -y " + s
}

// * * *

// NewSystem returns the interface to handle the given system.
func NewSystem(s system) pkg {
	switch s {
	case Debian:
		return new(debian)
	}
	panic("unreachable")
}
