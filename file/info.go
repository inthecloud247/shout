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

import "os"

// info represents a wrapper about os.FileInfo to append some functions.
type info struct{ fi os.FileInfo }

// NewInfo returns a info describing the named file.
func NewInfo(name string) (*info, error) {
	i, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return &info{i}, nil
}

// IsDir reports whether if it is a directory.
func (i *info) IsDir() bool {
	return i.fi.IsDir()
}

// IsFile reports whether it is a regular file.
func (i *info) IsFile() bool {
	return i.fi.Mode()&os.ModeType == 0
}

// IsModer reports whether it has read permission for the user.
func (i *info) IsModer() bool {
	return i.fi.Mode()&0400 != 0
}

// IsModew reports whether it has write permission for the user.
func (i *info) IsModew() bool {
	return i.fi.Mode()&0200 != 0
}

// IsModex reports whether it has execution permission for the user.
func (i *info) IsModex() bool {
	return i.fi.Mode()&0100 != 0
}

// IsMode reports whether it has the given permission.
// TODO: Can not be checked against 7.
func (i *info) IsMode(perm os.FileMode) bool {
	return i.fi.Mode()&perm != 0
}

// * * *

// IsDir reports whether if the named file is a directory.
func IsDir(name string) (bool, error) {
	i, err := NewInfo(name)
	if err != nil {
		return false, err
	}
	return i.IsDir(), nil
}

// IsFile reports whether the named file is a regular file.
func IsFile(name string) (bool, error) {
	i, err := NewInfo(name)
	if err != nil {
		return false, err
	}
	return i.IsFile(), nil
}

// IsModer reports whether the named file has read permission for the user.
func IsModer(name string) (bool, error) {
	i, err := NewInfo(name)
	if err != nil {
		return false, err
	}
	return i.IsModer(), nil
}

// IsModew reports whether the named file has write permission for the user.
func IsModew(name string) (bool, error) {
	i, err := NewInfo(name)
	if err != nil {
		return false, err
	}
	return i.IsModew(), nil
}

// IsModex reports whether the named file has execution permission for the user.
func IsModex(name string) (bool, error) {
	i, err := NewInfo(name)
	if err != nil {
		return false, err
	}
	return i.IsModex(), nil
}

// IsMode reports whether the named file has the permission perm.
func IsMode(name string, perm os.FileMode) (bool, error) {
	i, err := NewInfo(name)
	if err != nil {
		return false, err
	}
	return i.IsMode(perm), nil
}
