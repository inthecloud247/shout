// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package packager handles basic operations in package management systems.
//
// Important
//
// If you are going to use a package manager different to Deb, then you should
// check the options since I cann't test all.
//
// TODO
//
// Add managers of BSD systems.
//
// Use flag to do not show questions.
package packager

import (
	"errors"
	"os/exec"
)

type Packager interface {
	// Install installs a program.
	Install(...string) error

	// Remove removes a program.
	Remove(bool, ...string) error

	// Purge removes a program and its config files.
	Purge(bool, ...string) error

	// Clean erases downloaded archive files.
	Clean() error

	// Upgrade upgrades all the packages on the system.
	Upgrade() error
}

// * * *

// PackageType represents a package management system.
type PackageType int

const (
	Deb PackageType = iota + 1
	RPM
	Pacman
	Ebuild
	ZYpp
)

// New returns the interface to handle the package manager.
func New(p PackageType) Packager {
	switch p {
	case Deb:
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

// * * *

type packagerInfo struct {
	typ PackageType
	pkg Packager
}

// execPackagers is a list of executables of package managers.
var execPackagers = map[string]packagerInfo{
	"apt-get": packagerInfo{Deb, new(deb)},
	"yum":     packagerInfo{RPM, new(rpm)},
	"pacman":  packagerInfo{Pacman, new(pacman)},
	"emerge":  packagerInfo{Ebuild, new(ebuild)},
	"zypper":  packagerInfo{ZYpp, new(zypp)},
}

// Detect tries to get the package manager used in the system, looking for
// executables in directory "/usr/bin".
func Detect() (PackageType, Packager, error) {
	for k, v := range execPackagers {
		_, err := exec.LookPath("/usr/bin/" + k)
		if err == nil {
			return v.typ, v.pkg, nil
		}
	}
	return 0, nil, errors.New("package manager not found in directory /usr/bin")
}

// * * *

// runc executes a command logging its output if there is not any error.
func run(cmd string, arg ...string) error {
	_, err := exec.Command(cmd, arg...).CombinedOutput()
	if err != nil {
		return err
	}

	// log.Print(string(out)) // DEBUG
	return nil
}

// * * *

type packageSystem struct {
	isFirstInstall bool
}

// == Deb

type deb packageSystem

func (p deb) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/apt-get", "update"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	arg := []string{"install", "-y"}
	arg = append(arg, name...)
	return run("/usr/bin/apt-get", arg...)
}

func (deb) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"remove", "-y"}
	arg = append(arg, name...)
	if err := run("/usr/bin/apt-get", arg...); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/apt-get", "autoremove", "-y")
	}
	return nil
}

func (deb) Purge(isMetapackage bool, name ...string) error {
	arg := []string{"purge", "-y"}
	arg = append(arg, name...)
	if err := run("/usr/bin/apt-get", arg...); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/apt-get", "autoremove", "--purge", "-y")
	}
	return nil
}

func (deb) Clean() error {
	return run("/usr/bin/apt-get", "clean")
}

func (deb) Upgrade() error {
	if err := run("/usr/bin/apt-get", "update"); err != nil {
		return err
	}
	return run("/usr/bin/apt-get", "upgrade")
}

// http://fedoraproject.org/wiki/FAQ#How_do_I_install_new_software_on_Fedora.3F_Is_there_anything_like_APT.3F
// http://yum.baseurl.org/wiki/YumCommands
//
// == RPM

type rpm packageSystem

func (p rpm) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/yum", "update"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	arg := []string{"install"}
	arg = append(arg, name...)
	return run("/usr/bin/yum", arg...)
}

func (rpm) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"remove"}
	arg = append(arg, name...)
	return run("/usr/bin/yum", arg...)
}

func (rpm) Purge(isMetapackage bool, name ...string) error {
	return nil
}

func (rpm) Clean() error {
	return run("/usr/bin/yum", "clean", "packages")
}

func (rpm) Upgrade() error {
	return run("/usr/bin/yum", "update")
}

// https://wiki.archlinux.org/index.php/Pacman#Usage
// http://www.archlinux.org/pacman/pacman.8.html
//
// == Pacman

// "--noconfirm" bypasses the "Are you sure?" checks

type pacman packageSystem

func (p pacman) Install(name ...string) error {
	if p.isFirstInstall {
		p.isFirstInstall = false
		arg := []string{"-Syu", "--needed", "--noprogressbar"}
		arg = append(arg, name...)
		return run("/usr/bin/pacman", arg...)
	}

	arg := []string{"-S", "--needed", "--noprogressbar"}
	arg = append(arg, name...)
	return run("/usr/bin/pacman", arg...)
}

func (pacman) Remove(isMetapackage bool, name ...string) error {
	if isMetapackage {
		arg := []string{"-Rs"}
		arg = append(arg, name...)
		return run("/usr/bin/pacman", arg...)
	}

	arg := []string{"-R"}
	arg = append(arg, name...)
	return run("/usr/bin/pacman", arg...)
}

func (pacman) Purge(isMetapackage bool, name ...string) error {
	if isMetapackage {
		arg := []string{"-Rsn"}
		arg = append(arg, name...)
		return run("/usr/bin/pacman", arg...)
	}

	arg := []string{"-Rn"}
	arg = append(arg, name...)
	return run("/usr/bin/pacman", arg...)
}

func (pacman) Clean() error {
	return nil
}

func (pacman) Upgrade() error {
	return run("/usr/bin/pacman", "-Syu")
}

// http://www.gentoo.org/doc/en/handbook/handbook-x86.xml?part=2&chap=1
// http://www.gentoo-wiki.info/MAN_emerge
//
// == Ebuild

// --ask

type ebuild packageSystem

func (p ebuild) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/emerge", "--sync"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}
	return run("/usr/bin/emerge", name...)
}

func (ebuild) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"--unmerge"}
	arg = append(arg, name...)
	if err := run("/usr/bin/emerge", arg...); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/emerge", "--depclean")
	}
	return nil
}

func (ebuild) Purge(isMetapackage bool, name ...string) error {
	return nil
}

func (ebuild) Clean() error {
	return nil
}

func (ebuild) Upgrade() error {
	if err := run("/usr/bin/emerge", "--sync"); err != nil {
		return err
	}
	return run("/usr/bin/emerge", "--update", "--deep", "--with-bdeps=y", "--newuse world")
}

// http://en.opensuse.org/SDB:Zypper_usage
// http://www.openss7.org/man2html?zypper%288%29
//
// == ZYpp

// --no-confirm

type zypp packageSystem

func (p zypp) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/zypper", "refresh"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	arg := []string{"install", "--auto-agree-with-licenses"}
	arg = append(arg, name...)
	return run("/usr/bin/zypper", arg...)
}

func (zypp) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"remove"}
	arg = append(arg, name...)
	return run("/usr/bin/zypper", arg...)
}

func (zypp) Purge(isMetapackage bool, name ...string) error {
	return nil
}

func (zypp) Clean() error {
	return run("/usr/bin/zypper", "clean")
}

func (zypp) Upgrade() error {
	if err := run("/usr/bin/zypper", "refresh"); err != nil {
		return err
	}
	return run("/usr/bin/zypper", "up", "--auto-agree-with-licenses")
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
