errors [![Build Status](https://secure.travis-ci.org/101loops/config.png)](https://travis-ci.org/101loops/config)
======

This Go package adds an interface to deal with configuration files.

Noteworthy features:
- simple file syntax based on sections
- extraction of Duration and []byte
- variable substitution
- overwriting values by comman-line args and environment properties
- named sub-sections
- merging of multiple files


### Installation
`go get github.com/101loops/config`

### Documentation
[godoc.org](http://godoc.org/github.com/101loops/config)

### Credit
The package creates functionality on top of https://bitbucket.org/gosimple/conf/

### License
MIT (see LICENSE).