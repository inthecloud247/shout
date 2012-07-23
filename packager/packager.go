// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package packager handles basic operations in package management systems.
//
//
// Important
//
// If you are going to use a package manager different to DEB, then you should
// check the options since I cann't test all.
//
//
// TODO
// Add managers of BSD systems.  
// Use flag to do not show questions.
package packager

import (
	"os/exec"

	"github.com/kless/shout"
)

type Packager interface {
	// Install installs a program.
	Install(string) error

	// Remove removes a program.
	Remove(string, bool) error

	// Purge removes a program and its config files.
	Purge(string, bool) error

	// Clean erases downloaded archive files.
	Clean() error

	// Upgrade upgrades all the packages on the system.
	Upgrade() error
}

// PackageType represents a package management system.
type PackageType int

const (
	DEB PackageType = iota + 1
	RPM
	Pacman
	Ebuild
	ZYpp
)

// New returns the interface to handle the package manager.
func New(p PackageType) Packager {
	switch p {
	case DEB:
		return new(deb)
	case RPM:
		return new(rpm)
	case Pacman:
		return new(pacman)
	case Ebuild:
		return new(ebuild)
	case ZYpp:
		return new(zypp)
	}
	panic("unreachable")
}

// execPackagers is a list of executables of package managers.
var execPackagers = map[string]Packager{
	"apt-get": new(deb),
	"yum":     new(rpm),
	"pacman":  new(pacman),
	"emerge":  new(ebuild),
	"zypper":  new(zypp),
}

// Detect tries to get the package manager used in the system looking for
// executables in directory "/usr/bin".
func Detect() (p Packager, found bool) {
	for k, v := range execPackagers {
		_, err := exec.LookPath("/usr/bin/" + k)
		if err == nil {
			return v, true
		}
	}
	return nil, false
}

type packageSystem struct {
	isFirstInstall bool
}

// == DEB

// TODO: remove -s since it's to simulate (during testing)

type deb packageSystem

func (p deb) Install(name string) (err error) {
	if p.isFirstInstall {
		_, _, err = shout.Run("/usr/bin/apt-get update")
		p.isFirstInstall = false
	}
	_, _, err = shout.Run("/usr/bin/apt-get install -y -s " + name)
	return
}

func (deb) Remove(name string, isMetapackage bool) (err error) {
	_, _, err = shout.Run("/usr/bin/apt-get remove -y -s " + name)

	if isMetapackage && err == nil {
		_, _, err = shout.Run("/usr/bin/apt-get autoremove -y -s")
	}
	return
}

func (deb) Purge(name string, isMetapackage bool) (err error) {
	_, _, err = shout.Run("/usr/bin/apt-get purge -y -s " + name)

	if isMetapackage && err == nil {
		_, _, err = shout.Run("/usr/bin/apt-get autoremove --purge -y -s")
	}
	return
}

func (deb) Clean() (err error) {
	_, _, err = shout.Run("/usr/bin/apt-get clean")
	return
}

func (deb) Upgrade() (err error) {
	_, _, err = shout.Run("/usr/bin/apt-get update")
	_, _, err = shout.Run("/usr/bin/apt-get upgrade")
	return
}

// http://fedoraproject.org/wiki/FAQ#How_do_I_install_new_software_on_Fedora.3F_Is_there_anything_like_APT.3F
// http://yum.baseurl.org/wiki/YumCommands
//
// == RPM

type rpm packageSystem

func (p rpm) Install(name string) (err error) {
	if p.isFirstInstall {
		_, _, err = shout.Run("/usr/bin/yum update")
		p.isFirstInstall = false
	}
	_, _, err = shout.Run("/usr/bin/yum install " + name)
	return
}

func (rpm) Remove(name string, isMetapackage bool) (err error) {
	_, _, err = shout.Run("/usr/bin/yum remove " + name)
	return
}

func (rpm) Purge(name string, isMetapackage bool) (err error) {
	return nil
}

func (rpm) Clean() (err error) {
	_, _, err = shout.Run("/usr/bin/yum clean packages")
	return
}

func (rpm) Upgrade() (err error) {
	_, _, err = shout.Run("/usr/bin/yum update")
	return
}

// https://wiki.archlinux.org/index.php/Pacman#Usage
// http://www.archlinux.org/pacman/pacman.8.html
//
// == Pacman

// "--noconfirm" bypasses the "Are you sure?" checks

type pacman packageSystem

func (p pacman) Install(name string) (err error) {
	if p.isFirstInstall {
		_, _, err = shout.Run("/usr/bin/pacman -Syu --needed --noprogressbar " + name)
		p.isFirstInstall = false
	} else {
		_, _, err = shout.Run("/usr/bin/pacman -S --needed  --noprogressbar " + name)
	}
	return
}

func (pacman) Remove(name string, isMetapackage bool) (err error) {
	if isMetapackage {
		_, _, err = shout.Run("/usr/bin/pacman -Rs " + name)
	} else {
		_, _, err = shout.Run("/usr/bin/pacman -R " + name)
	}
	return
}

func (pacman) Purge(name string, isMetapackage bool) (err error) {
	if isMetapackage {
		_, _, err = shout.Run("/usr/bin/pacman -Rsn " + name)
	} else {
		_, _, err = shout.Run("/usr/bin/pacman -Rn " + name)
	}
	return
}

func (pacman) Clean() (err error) {
	return nil
}

func (pacman) Upgrade() (err error) {
	_, _, err = shout.Run("/usr/bin/pacman -Syu")
	return
}

// http://www.gentoo.org/doc/en/handbook/handbook-x86.xml?part=2&chap=1
// http://www.gentoo-wiki.info/MAN_emerge
//
// == Ebuild

// --ask

type ebuild packageSystem

func (p ebuild) Install(name string) (err error) {
	if p.isFirstInstall {
		_, _, err = shout.Run("/usr/bin/emerge --sync")
		p.isFirstInstall = false
	}
	_, _, err = shout.Run("/usr/bin/emerge " + name)
	return
}

func (ebuild) Remove(name string, isMetapackage bool) (err error) {
	_, _, err = shout.Run("/usr/bin/emerge --unmerge " + name)

	if isMetapackage && err == nil {
		_, _, err = shout.Run("/usr/bin/emerge --depclean")
	}
	return
}

func (ebuild) Purge(name string, isMetapackage bool) (err error) {
	return nil
}

func (ebuild) Clean() (err error) {
	return nil
}

func (ebuild) Upgrade() (err error) {
	_, _, err = shout.Run("/usr/bin/emerge --sync")
	_, _, err = shout.Run("/usr/bin/emerge --update --deep --with-bdeps=y --newuse world")
	return
}

// http://en.opensuse.org/SDB:Zypper_usage
// http://www.openss7.org/man2html?zypper%288%29
//
// == ZYpp

// --no-confirm

type zypp packageSystem

func (p zypp) Install(name string) (err error) {
	if p.isFirstInstall {
		_, _, err = shout.Run("/usr/bin/zypper refresh")
		p.isFirstInstall = false
	}
	_, _, err = shout.Run("/usr/bin/zypper install --auto-agree-with-licenses " + name)
	return
}

func (zypp) Remove(name string, isMetapackage bool) (err error) {
	_, _, err = shout.Run("/usr/bin/zypper remove " + name)
	return
}

func (zypp) Purge(name string, isMetapackage bool) (err error) {
	return nil
}

func (zypp) Clean() (err error) {
	_, _, err = shout.Run("/usr/bin/zypper clean")
	return
}

func (zypp) Upgrade() (err error) {
	_, _, err = shout.Run("/usr/bin/zypper refresh")
	_, _, err = shout.Run("/usr/bin/zypper up --auto-agree-with-licenses")
	return
}

/* TODO: maybe be needed ahead
type system int

// System distributions.
const (
	Arch system = iota + 1
	CentOS
	Debian
	Fedora
	Gentoo
	Mandriva
	OpenSUSE
	RHEL
	Scientific
	Slackware
	SUSE
	Turbolinux
	Ubuntu
)

var ListSystems = map[system]string{
	Arch:       "Arch",
	CentOS:     "CentOS",
	Debian:     "Debian",
	Fedora:     "Fedora",
	Gentoo:     "Gentoo",
	Mandriva:   "Mandriva",
	OpenSUSE:   "openSUSE",
	RHEL:       "Red Hat Enterprise Linux",
	Scientific: "Scientific Linux",
	Slackware:  "Slackware",
	SUSE:       "SUSE Enterprise Linux",
	Turbolinux: "Turbolinux",
	Ubuntu:     "Ubuntu",
}*/
