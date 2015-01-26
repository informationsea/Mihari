Mihari
======

Construct data analysis pipeline without writing shell script or makefile.

How to use?
-----------

1. Install `mihari` to somewhere
2. Run `mihari --init` on your working directory.
3. Run your analysis commands with `mihari ` prefix.
4. Run `mihari --makefile`. mihari generates `Makefile`
   automatically.
5. Run `make` to reproduce generated files.
6. Run `make clean` to remove generated files.

Warning
-------

This software has a lot of bugs. This software cannot track all file
operations currently, and a format of log is not stable. PLEASE RUN
`make -n` AND CHECK WHAT MIHARI WILL DO!

How to build?
-------------

1. Install [GO language](https://golang.org/dl/) 1.4 or later
2. Set up `$GOPATH`
3. Install Go-bindata `go get -u github.com/jteeuwen/go-bindata/...`
4. Set up `export PATH=$GOPATH/bin:$PATH`
5. Run `make`

License
-------

    Miahri : Construct data analysis pipeline without writing
    Copyright (C) 2015 OKAMURA Yasunobu

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

