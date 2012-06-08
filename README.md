shout
=====

== The beginning

The shell scripting gets to be no-maintainable code; its sintaxis is very
cryptic and it's very hard to debug. In addition, these negative points increase
at the same time as the size of the program grows.

Here is where [Go][] comes.

Go is a compiled, garbage-collected, concurrent programming language that
provides the efficiency of a statically typed compiled language with the ease of
programming of a dynamic language.

The compilers can target the FreeBSD, Linux, NetBSD, OpenBSD, Mac OS X (Snow
Leopard/Lion), and Windows operating systems and the 32-bit (386) and 64-bit
(amd64) x86 processor architectures, and 32-bit ARM for Linux.

It has a simple and clean [sintaxis][], an error handling that does code more
reliable, and a full [standard library][] although the package [os][] will be
the main one to use in system scripts.  
Whatever administrator without great knowledge about programming can built basic
scripts fastly with the help of this [tutorial for novices][].

The advantage of a shell script, versus a compiled program, is that allows an
easy modification and locating of sources. But Go also can do the same using
[goplay][].

[Go]: http://golang.org/
[sintaxis]: http://golang.org/ref/spec
[standard library]: http://golang.org/pkg/
[os]: http://golang.org/pkg/os/
[tutorial for novices]: http://go-book.appspot.com/
[goplay]: https://github.com/kless/goplay

== External commands

The main tool of this package is the function *Run* which lets to run system
commands under a new process. It handles pipes, environment variables, and does
pattern expansion just as in the Bash shell.

== Editing

The editing of files is very important in the shell scripting to working with
the configuration files. shout has a great number of functions related to it,
avoiding to have to use an external command to get the same result, and with the
advantage of that it is created automatically a backup before of editing a file.


## Installation

	go get github.com/kless/shout

To run the tests:

	cd ${GOPATH//:*}/src/github.com/kless/shout && go test && cd -


## Configuration

*NewEdit* creates a new struct *edit* which has a variable, *CommentChar*,
with a value by default, '#'. That value is the character used in comments.


## Copyright and licensing

*Copyright 2012  The "shout" Authors*. See file AUTHORS and CONTRIBUTORS.  
Unless otherwise noted, the source files are distributed under the
*Apache License, version 2.0* found in the LICENSE file.


* * *
*Generated by [gowizard](https://github.com/kless/gowizard)*

