// Copyright 2012  The "Shout" Authors
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
	"log"
	"os"
	"os/exec"
)

var (
	_log    *log.Logger
	logFile *os.File
)

// TODO: see /var/log/apt/history.log to have a similar log schema.
func init() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")

	if os.Getuid() != 0 {
		log.Fatal("you have to be root")
	}

	f, err := os.OpenFile("/var/log/shout/packager.log", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Seek(0, os.SEEK_END)
	if err != nil {
		log.Fatal(err)
	}

	logFile = f
	_log = log.New(logFile, "", log.LstdFlags)
}

// CloseLogfile closes the log file.
func CloseLogfile() error {
	return logFile.Close()
}

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

type packagerInfo struct {
	pkg Packager
	typ PackageType
}

// execPackagers is a list of executables of package managers.
var execPackagers = map[string]packagerInfo{
	"apt-get": packagerInfo{new(deb), Deb},
	"yum":     packagerInfo{new(rpm), RPM},
	"pacman":  packagerInfo{new(pacman), Pacman},
	"emerge":  packagerInfo{new(ebuild), Ebuild},
	"zypper":  packagerInfo{new(zypp), ZYpp},
}

// Detect tries to get the package manager used in the system, looking for
// executables in directory "/usr/bin".
func Detect() (Packager, PackageType, error) {
	for k, v := range execPackagers {
		_, err := exec.LookPath("/usr/bin/" + k)
		if err == nil {
			return v.pkg, v.typ, nil
		}
	}
	return nil, 0, errors.New("package manager not found in directory /usr/bin")
}

// runc executes a command logging its output if there is not any error.
func run(name string, arg ...string) error {
	out, err := exec.Command(name, arg...).CombinedOutput()
	if err != nil {
		return err
	}
	_log.Print(out)
	return nil
}

type packageSystem struct {
	isFirstInstall bool
}

// == Deb

type deb packageSystem

func (p deb) Install(name string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/apt-get", "update"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	return run("/usr/bin/apt-get", "install", "-y", name)
}

func (deb) Remove(name string, isMetapackage bool) error {
	if err := run("/usr/bin/apt-get", "remove", "-y", name); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/apt-get", "autoremove", "-y")
	}
	return nil
}

func (deb) Purge(name string, isMetapackage bool) error {
	if err := run("/usr/bin/apt-get", "purge", "-y", name); err != nil {
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

func (p rpm) Install(name string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/yum", "update"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	return run("/usr/bin/yum", "install", name)
}

func (rpm) Remove(name string, isMetapackage bool) error {
	return run("/usr/bin/yum", "remove", name)
}

func (rpm) Purge(name string, isMetapackage bool) error {
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

func (p pacman) Install(name string) error {
	if p.isFirstInstall {
		p.isFirstInstall = false
		return run("/usr/bin/pacman", "-Syu", "--needed", "--noprogressbar", name)
	}
	return run("/usr/bin/pacman", "-S", "--needed", "--noprogressbar", name)
}

func (pacman) Remove(name string, isMetapackage bool) error {
	if isMetapackage {
		return run("/usr/bin/pacman", "-Rs", name)
	}
	return run("/usr/bin/pacman", "-R", name)
}

func (pacman) Purge(name string, isMetapackage bool) error {
	if isMetapackage {
		return run("/usr/bin/pacman", "-Rsn", name)
	}
	return run("/usr/bin/pacman", "-Rn", name)
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

func (p ebuild) Install(name string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/emerge", "--sync"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}
	return run("/usr/bin/emerge", name)
}

func (ebuild) Remove(name string, isMetapackage bool) error {
	if err := run("/usr/bin/emerge", "--unmerge", name); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/emerge", "--depclean")
	}
	return nil
}

func (ebuild) Purge(name string, isMetapackage bool) error {
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

func (p zypp) Install(name string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/zypper", "refresh"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}
	return run("/usr/bin/zypper", "install", "--auto-agree-with-licenses", name)
}

func (zypp) Remove(name string, isMetapackage bool) error {
	return run("/usr/bin/zypper", "remove", name)
}

func (zypp) Purge(name string, isMetapackage bool) error {
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
