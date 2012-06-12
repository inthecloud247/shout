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
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const _BACKUP_SUFFIX = "+[1-9]~" // suffix pattern added to backup's file name

// Backup creates a backup of the named file.
//
// The schema used for the new name is: {name}\+[1-9]~
//   name: The original file name.
//   + : Character used to separate the file name from rest.
//   number: A number from 1 to 9, using rotation.
//   ~ : To indicate that it is a backup, just like it is used in Unix systems.
func Backup(name string) error {
	// Check if it is empty
	info, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.Size() == 0 {
		return nil
	}

	files, err := filepath.Glob(name + _BACKUP_SUFFIX)
	if err != nil {
		return err
	}

	// Number rotation
	numBackup := byte(1)

	if len(files) != 0 {
		lastFile := files[len(files)-1]
		numBackup = lastFile[len(lastFile)-2] + 1 // next number

		if numBackup > '9' {
			numBackup = '1'
		}
	} else {
		numBackup = '1'
	}

	_, err = Copy(name, fmt.Sprintf("%s+%s~", name, string(numBackup)))
	return err
}

// Copy copies file in source to file in dest preserving the mode attributes.
func Copy(source, dest string) (int64, error) {
	// Don't backup files of backup.
	if dest[len(dest)-1] != '~' {
		if err := Backup(dest); err != nil {
			return 0, err
		}
	}

	srcFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	srcInfo, err := os.Stat(source)
	if err != nil {
		return 0, err
	}

	dstFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

// Create creates a new file with b bytes.
func Create(name string, b []byte) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(b)
	return err
}

// CreateString is like Create, but writes the contents of string s rather than
// an array of bytes.
func CreateString(name, s string) error {
	return Create(name, []byte(s))
}

// Overwrite truncates the named file to zero and writes len(b) bytes. It
// returns an error, if any.
func Overwrite(name string, b []byte) error {
	if err := Backup(name); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}

// OverwriteString is like Overwrite, but writes the contents of string s rather
// than an array of bytes.
func OverwriteString(name, s string) error {
	return Overwrite(name, []byte(s))
}
