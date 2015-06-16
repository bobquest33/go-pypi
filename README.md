# go-pypi

[View this on GitHub Pages](http://ccpgames.github.io/go-pypi/)

[![Build status](https://img.shields.io/travis/ccpgames/go-pypi.svg)](https://travis-ci.org/ccpgames/go-pypi)
[![Coverage Status](https://img.shields.io/coveralls/ccpgames/go-pypi.svg)](https://coveralls.io/r/ccpgames/go-pypi?branch=master)
[![License](https://img.shields.io/github/license/ccpgames/go-pypi.svg)](https://github.com/ccpgames/go-pypi/blob/master/LICENSE)

Using PyPI with Go

## install

```bash
$ go get github.com/ccpgames/go-pypi
```

## use

Get the latest version of packages by default:

```bash
$ go-pypi requests
requests-2.7.0-py2.py3-none-any.whl downloaded (470641 bytes)
requests-2.7.0.tar.gz downloaded (451723 bytes)
```

Get a specific version (and/or multiple packages at once):

```bash
$ go-pypi requests=1.0.0 Flask
requests-1.0.0.tar.gz downloaded (335548 bytes)
Flask-0.10.1.tar.gz downloaded (544247 bytes)
```

Get a specific format:

```bash
$ go-pypi --extension=whl requests
requests-2.7.0-py2.py3-none-any.whl downloaded (470641 bytes)
```

Get from a different PyPI server:

```bash
$ go-pypi --url=https://pypi.yourcompany.com/pypi requests
requests-2.7.0-py2.py3-none-any.whl downloaded (470641 bytes)
requests-2.7.0.tar.gz downloaded (451723 bytes)
```
