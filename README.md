[![Build status](https://secure.travis-ci.org/ostrost/ostent.png?branch=master)](https://travis-ci.org/ostrost/ostent)
[![Sourcegraph](https://sourcegraph.com/api/repos/github.com/ostrost/ostent/.badges/status.png)](https://sourcegraph.com/github.com/ostrost/ostent)
[![GoDoc](https://godoc.org/github.com/ostrost/ostent?status.svg)](https://godoc.org/github.com/ostrost/ostent)

`ostent` displays current system metrics. [**Demo** here](http://demo.ostrost.com/)

![Screenshot](https://www.ostrost.com/ostent/screenshot.png)

Install & run with `curl -sSL https://github.com/ostrost/ostent/raw/master/ostent.sh | sh`

It's a single executable without dependecies. Once installed,
it will self-upgrade whenever there's new release.

Platforms
---------

   - Linux [64-bit](https://github.com/ostrost/ostent/releases/download/v0.1.9/Linux.x86_64) | [32-bit](https://github.com/ostrost/ostent/releases/download/v0.1.9/Linux.i686)
   - [Darwin](https://github.com/ostrost/ostent/releases/download/v0.1.9/Darwin.x86_64)
   - _Expect \*BSD builds surely_

Binaries distributed by [GitHub Releases](https://github.com/ostrost/ostent/releases)

Usage
-----

`ostent` accepts optional `-bind` argument to set specific IP and/or
port to bind to, otherwise any machine IP and port 8050 by default.

   - `ostent -bind 127.1` # [http://127.0.0.1:8050/](http://127.0.0.1:8050/)
   - `ostent -bind 192.168.1.10:8051` # port 8051
   - `ostent -bind 8052` # any IP, port 8052

`-update` sets collection interval (1 second by default),
append `s` for seconds, `m` for minutes: `5s`, `1m` etc.

Run it, it'll give the link(s) to open in a browser.

Running the code
----------------

Have your GOPATH environment set,
[gvm](https://github.com/moovweb/gvm) is a must.

1. **`go get github.com/ostrost/ostent/ostent`**
2. `ostent` to run.

For rebuilding the code and assets:

1. Find `src/github.com/ostrost/ostent` directory in GOPATH.
2. Run `make init` once and later for packages update (think `go get -u`)
3. `make` or `make al` when `rerun` does rebuilding.

Repeat 3. every time sources (esp. assets) change.
[rerun](https://github.com/skelterjohn/rerun) does live-reloading run:
`rerun github.com/ostrost/ostent/ostent`

**For a fork**, to preserve import paths and packages namespace,
clone your fork as if it was `github.com/ostrost/ostent` package for Go:

1. `go get github.com/ostrost/ostent/ostent`
2. Find `src/github.com/ostrost/ostent` directory in GOPATH.
3. Replace it with you fork clone.
4. Continue with rebuilding steps above.

Make
----

`make` rebuilds these **commited to the repo** files:
- `src/share/templates/bindata.*.go`
- `src/share/assets/bindata.*.go`
- `src/share/assets/js/devel/milk/*.js`
- `src/share/assets/js/devel/gen/*.js`
- `src/share/templates/*.html`
- `src/share/assets/css/*.css`
- `src/share/tmp/jsassets.d`
- `src/share/tmp/*.jsx`

If you don't change source files, content re-generated
should not differ from the commited. Whenever
src/share/{amber.templates,style,coffee} modified,
you have to re-make.

Additional tools required for assets rebuilding:
- [Sass](http://sass-lang.com/install)
- [react-tools](https://www.npmjs.org/package/react-tools)
- [uglify-js](https://www.npmjs.org/package/uglify-js)

Go packages
-----------

`./ostent` is the main (_as in [Go Program execution](http://golang.org/ref/spec#Program_execution)_) package:
rerun will find `main.devel.go` file; the other `main.production.go`
(used when building with `-tags production`) is the init code for
the distributed binaries: also includes
[goagain](https://github.com/rcrowley/goagain) recovering and
self-upgrading via [go-update](https://github.com/inconshreveable/go-update).

`src/amberp/amberpp` is templates compiler, used with make.

The assets
----------

The binaries, to be stand-alone, have the assets and templates embeded.
Unless you specifically build with `-tags production` (e.g with make),
the content is not embeded for the ease of development:
with `rerun`, asset requests are served from the actual files.
Production-built `ostent extract-assets` can be used to copy (extract) assets on disk.
