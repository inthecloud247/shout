// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Installer manages files related to the system.
//
// During the development, it could be used the command:
//
// sudo env PATH=$PATH GOPATH=$GOPATH go run install.go <flag...>
package main

import (
	"flag"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/kless/sysuser"
)

var (
	logDir   = "/var/log/shout"
	logFiles = []string{"shout"}
)

func Install() {
	err := os.Mkdir(logDir, 0755)
	if err != nil {
		if os.IsExist(err) {
			return
		}
		log.Fatal(err)
	}

	// TODO: handle group in Windows
	group, err := sysuser.LookupGroup("adm")
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range logFiles {
		f, err := os.OpenFile(path.Join(logDir, name)+".log", os.O_CREATE, 0640)
		if err != nil {
			log.Fatal(err)
		}

		f.Chown(0, group.Gid)
		f.Close()
	}
}

func Remove() {
	// Remove binary files.
}

func Purge() {
	Remove()

	err := os.RemoveAll(logDir)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
}

// * * *

func init() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")

	if runtime.GOOS == "windows" {
		log.Fatal("TODO: add a directory to save logs in Windows")
	}
	if os.Getuid() != 0 {
		log.Fatal("you have to be root")
	}
}

func main() {
	install := flag.Bool("i", false, "install")
	remove := flag.Bool("r", false, "remove")
	purge := flag.Bool("p", false, "purge")

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	if *purge {
		Purge()
	} else if *remove {
		Remove()
	}
	if *install {
		Install()
	}
}
