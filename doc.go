// Copyright 2012  The "Shout" Authors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package shout handles the shell scripting.

The main tool of this package is the function Run which lets to run system
commands under a new process. It handles pipes, environment variables, and does
pattern expansion just as in the Bash shell.

The editing of files is very important in the shell scripting to working with
the configuration files. Shout has a great number of functions related to it,
avoiding to have to use an external command to get the same result, and with the
advantage of that it is created automatically a backup before of editing a file.


Configuration

NewEdit creates a new struct, edit, which has a variable, CommentChar,
with a value by default, '#'. That value is the character used in comments.
*/
package shout
