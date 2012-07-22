// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package system handles the commands to install and remove programs into an
// operating system.
package system

type Packager interface {
	// Install returns the command to install a program.
	Install(s string) string

	// Remove returns the command to remove a program.
	Remove(s string) string

	// Purge returns the command to purge a program.
	Purge(s string) string
}

// system represents the operating system or Linux flavor.
type system int

// Operating system or Linux flavor.
const (
	Debian system = iota + 1
)

// NewSystem returns the interface to handle the given system.
func NewSystem(s system) Packager {
	switch s {
	case Debian:
		return new(debian)
	}
	panic("unreachable")
}

//
// == Systems

type debian system

func (d debian) Install(s string) string {
	return "apt-get install -y " + s
}

func (d debian) Remove(s string) string {
	return "apt-get remove -y " + s
}

func (d debian) Purge(s string) string {
	return "apt-get purge -y " + s
}
