// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package system handles basic operations in package managers.
package system

import "github.com/kless/shout"

type Packager interface {
	// Update retrieves new lists of packages.
	Update() error

	// Clean erases downloaded archive files.
	Clean() error

	// Install runs the command to install a program.
	Install(string) error

	// Remove runs the command to remove a program.
	Remove(string, bool) error

	// Purge runs the command to remove a program and its config files.
	Purge(string, bool) error
}

// packageSystem represents a package management system.
type packageSystem int

// Operating system or Linux flavor.
const (
	Debian packageSystem = iota + 1
)

var ListSystems = map[packageSystem]string{
	Debian: "Debian",
}

// NewSystem returns the interface to handle the given system.
func NewSystem(s packageSystem) Packager {
	switch s {
	case Debian:
		return new(apt)
	}
	panic("unreachable")
}

// Packages represents a package name in every system.
//type Packages map[packageSystem]string

//
// == Systems

// == APT

// TODO: remove -s since it's to simulate (during testing)

type apt packageSystem

func (apt) Update() error {
	_, _, err := shout.Run("/usr/bin/apt-get update")
	return err
}

func (apt) Install(name string) (err error) {
	_, _, err = shout.Run("/usr/bin/apt-get install -y -s " + name)
	return
}

func (apt) Remove(name string, isMetapackage bool) (err error) {
	_, _, err = shout.Run("/usr/bin/apt-get remove -y -s " + name)

	if isMetapackage && err != nil {
		_, _, err = shout.Run("/usr/bin/apt-get autoremove -y -s")
	}
	return
}

func (apt) Purge(name string, isMetapackage bool) (err error) {
	_, _, err = shout.Run("/usr/bin/apt-get purge -y -s " + name)

	if isMetapackage && err != nil {
		_, _, err = shout.Run("/usr/bin/apt-get autoremove --purge -y -s")
	}
	return
}

func (apt) Clean() error {
	_, _, err := shout.Run("/usr/bin/apt-get clean")
	return err
}

// == YUM

type yum packageSystem

func (yum) Update() error {
	_, _, err := shout.Run("yum update")
	return err
}

func (yum) Install(name string) (err error) {
	_, _, err = shout.Run("yum install " + name)
	return
}

func (yum) Remove(name string) (err error) {
	_, _, err = shout.Run("yum remove " + name)
	return
}

/*func (yum) Purge(name string) (err error) {
	_, _, err = shout.Run("yum remove " + name)
	return
}*/

func (yum) Clean() error {
	_, _, err := shout.Run("yum clean packages")
	return err
}

// SUSE, Gentoo, Mandriva, Slackware, Fedora, Turbolinux, Arch

