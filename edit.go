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

package shout

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

// editDefault represents the vaues by default to set in type edit.
type editDefault struct {
	CommentChar string // character used in comments
	//DoBackup    bool   // do backup before of edit?
}

// Values by default for type edit.
var _editDefault = editDefault{"#"}

// edit represents the file to edit.
type edit struct {
	editDefault
	file *os.File
	buf  *bufio.ReadWriter
}

type replacer struct {
	search, replace string
}

type replacerAtLine struct {
	line, search, replace string
}

// NewEdit opens a file to edit; it is created a backup.
func NewEdit(name string) (*edit, error) {
	if err := Backup(name); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return &edit{
		_editDefault,
		file,
		bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file)),
	}, nil
}

// Append writes len(b) bytes at the end of the File. It returns an error, if any.
func (e *edit) Append(b []byte) error {
	return e.write(b, os.SEEK_END)
}

// AppendString is like Append, but writes the contents of string s rather than
// an array of bytes.
func (e *edit) AppendString(s string) error {
	return e.write([]byte(s), os.SEEK_END)
}

// Close closes the file.
func (e *edit) Close() error {
	return e.file.Close()
}

// Comment inserts the comment character in lines that mach any regular expression
// in reLine.
func (e *edit) Comment(reLine []string) error {
	allReSearch := make([]*regexp.Regexp, len(reLine))
	char := []byte(e.CommentChar + " ")

	for _, v := range reLine {
		if re, err := regexp.Compile(v); err != nil {
			return err
		} else {
			allReSearch = append(allReSearch, re)
		}
	}

	newContent := new(bytes.Buffer)
	isNew := false

	// Check every line.
	for {
		line, err := e.buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		for _, v := range allReSearch {
			if v.Match(line) {
				line = append(char, line...)

				if !isNew {
					isNew = true
				}
			}
			if _, err = newContent.Write(line); err != nil {
				return err
			}
		}
	}

	if isNew {
		return e.rewrite(newContent.Bytes())
	}
	return nil
}

// CommentOut removes the comment character of lines that mach any regular expression
// in reLine.
func (e *edit) CommentOut(reLine []string) error {
	allSearch := make([]replacerAtLine, len(reLine))

	for i, v := range reLine {
		allSearch[i] = replacerAtLine{
			v,
			"[[:space:]]*" + e.CommentChar + "[[:space:]]*",
			"",
		}
	}

	return e.ReplaceAtLineN(allSearch, 1)
}

// Insert writes len(b) bytes at the start of the File. It returns an error, if any.
func (e *edit) Insert(b []byte) error {
	return e.write(b, os.SEEK_SET)
}

// InsertString is like Insert, but writes the contents of string s rather than
// an array of bytes.
func (e *edit) InsertString(s string) error {
	return e.write([]byte(s), os.SEEK_SET)
}

// Replace replaces all regular expressions mathed in r.
func (e *edit) Replace(r []replacer) error {
	return e.genReplace(r, -1)
}

// ReplaceN replaces regular expressions mathed in r. The count determines the
// number to match:
//   n > 0: at most n matches
//   n == 0: the result is none
//   n < 0: all matches
func (e *edit) ReplaceN(r []replacer, n int) error {
	return e.genReplace(r, n)
}

// ReplaceAtLine replaces all regular expressions mathed in r, if the line is
// matched at the first.
func (e *edit) ReplaceAtLine(r []replacerAtLine) error {
	return e.genReplaceAtLine(r, -1)
}

// ReplaceAtLine replaces regular expressions mathed in r, if the line is
// matched at the first. The count determines the
// number to match:
//   n > 0: at most n matches
//   n == 0: the result is none
//   n < 0: all matches
func (e *edit) ReplaceAtLineN(r []replacerAtLine, n int) error {
	return e.genReplaceAtLine(r, n)
}

// Generic Replace: replaces a number of regular expressions matched in r.
func (e *edit) genReplace(r []replacer, n int) error {
	if n == 0 {
		return nil
	}

	content, err := ioutil.ReadAll(e.buf)
	if err != nil {
		return err
	}

	var lastContent []byte
	isNew := false

	for _, v := range r {
		reSearch, err := regexp.Compile(v.search)
		if err != nil {
			return err
		}

		if n < 0 {
			lastContent = reSearch.ReplaceAllLiteral(content, []byte(v.replace))
		} else {
			repl := []byte(v.replace)
			lastContent = content

			for _, idx := range reSearch.FindAllSubmatchIndex(content, n) {
				lastContent = append(lastContent[:idx[0]], append(repl, lastContent[idx[1]:]...)...)
			}
		}

		if !bytes.Equal(content, lastContent) {
			content = lastContent
			if !isNew {
				isNew = true
			}
		}
	}

	if isNew {
		return e.rewrite(content)
	}
	return nil
}

// Generic ReplaceAtLine: replaces a number of regular expressions matched in r,
// if the line is matched at the first.
func (e *edit) genReplaceAtLine(r []replacerAtLine, n int) error {
	if n == 0 {
		return nil
	}

	// == Cache the regular expressions
	allReLine := make([]*regexp.Regexp, len(r))
	allReSearch := make([]*regexp.Regexp, len(r))
	allRepl := make([][]byte, len(r))

	for _, v := range r {
		if reLine, err := regexp.Compile(v.line); err != nil {
			return err
		} else {
			allReLine = append(allReLine, reLine)
		}

		if reSearch, err := regexp.Compile(v.search); err != nil {
			return err
		} else {
			allReSearch = append(allReSearch, reSearch)
		}

		allRepl = append(allRepl, []byte(v.replace))
	}

	newContent := new(bytes.Buffer)
	isNew := false

	// Replace every line, if it maches
	for {
		line, err := e.buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		newLine := line

		for i, _ := range r {
			if allReLine[i].Match(line) {
				if n < 0 {
					newLine = allReSearch[i].ReplaceAllLiteral(newLine, allRepl[i])
				} else {
					for _, idx := range allReSearch[i].FindAllSubmatchIndex(newLine, n) {
						newLine = append(newLine[:idx[0]],
							append(allRepl[i], newLine[idx[1]:]...)...)
					}
				}
			}

			if _, err = newContent.Write(newLine); err != nil {
				return err
			}
			if !isNew && !bytes.Equal(line, newLine) {
				isNew = true
			}
		}
	}

	if isNew {
		return e.rewrite(newContent.Bytes())
	}
	return nil
}

func (e *edit) rewrite(b []byte) error {
	if _, err := e.file.Seek(0, os.SEEK_SET); err != nil {
		return err
	}

	n, err := e.file.Write(b)
	if err != nil {
		return err
	}
	if err = e.file.Truncate(int64(n)); err != nil {
		return err
	}

	return e.file.Sync()
}

func (e *edit) write(b []byte, seek int) error {
	_, err := e.file.Seek(0, seek)
	if err != nil {
		return err
	}

	_, err = e.file.Write(b)
	return err
}

// * * *

// Append writes len(b) bytes at the end of the file filename. It returns an
// error, if any. The file is backed up.
func Append(filename string, b []byte) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.Append(b)
}

// AppendString is like Append, but writes the contents of string s rather than
// an array of bytes.
func AppendString(filename, s string) error {
	return Append(filename, []byte(s))
}

// Comment inserts the comment character in lines that mach the regular expression
// in reLine, in the file filename.
func Comment(filename, reLine string) error {
	return CommentM(filename, []string{reLine})
}

// CommentM inserts the comment character in lines that mach any regular expression
// in reLine, in the file filename.
func CommentM(filename string, reLine []string) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.Comment(reLine)
}

// CommentOut removes the comment character of lines that mach the regular expression
// in reLine, in the file filename.
func CommentOut(filename, reLine string) error {
	return CommentOutM(filename, []string{reLine})
}

// CommentOutM removes the comment character of lines that mach any regular expression
// in reLine, in the file filename.
func CommentOutM(filename string, reLine []string) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.CommentOut(reLine)
}

// Insert writes len(b) bytes at the start of the file filename. It returns an
// error, if any. The file is backed up.
func Insert(filename string, b []byte) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.Insert(b)
}

// InsertString is like Insert, but writes the contents of string s rather than
// an array of bytes.
func InsertString(filename, s string) error {
	return Insert(filename, []byte(s))
}

// Overwrite truncates the file filename to zero and writes len(b) bytes. It
// returns an error, if any.
func Overwrite(filename string, b []byte) error {
	if err := Backup(filename); err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}

// OverwriteString is like Overwrite, but writes the contents of string s rather
// than an array of bytes.
func OverwriteString(filename, s string) error {
	return Overwrite(filename, []byte(s))
}

// Replace replaces all regular expressions mathed in r for the file filename.
func Replace(filename string, r []replacer) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.genReplace(r, -1)
}

// ReplaceN replaces a number of regular expressions mathed in r for the file
// filename.
func ReplaceN(filename string, r []replacer, n int) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.genReplace(r, n)
}

// ReplaceAtLine replaces all regular expressions mathed in r for the file
// filename, if the line is matched at the first.
func ReplaceAtLine(filename string, r []replacerAtLine) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.genReplaceAtLine(r, -1)
}

// ReplaceAtLineN replaces a number of regular expressions mathed in r for the
// file filename, if the line is matched at the first.
func ReplaceAtLineN(filename string, r []replacerAtLine, n int) error {
	e, err := NewEdit(filename)
	if err != nil {
		return err
	}
	defer e.Close()

	return e.genReplaceAtLine(r, n)
}
