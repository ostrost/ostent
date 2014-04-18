`ostent` displays current system metrics. [**Demo** here](http://demo.ostrost.com/)

![screenshot](https://github.com/rzab/ostent/raw/master/screenshot.png)

Install & run with `curl -sSL https://github.com/rzab/ostent/raw/master/ostent.sh | sh`

It's a single executable without dependecies. Once installed,
it will self-upgrade whenever there's new release.

Platforms
---------

   - [Linux 64bits](https://github.com/rzab/ostent/releases/download/v0.1.3/Linux.x86_64)
   - [Linux 32bits](https://github.com/rzab/ostent/releases/download/v0.1.3/Linux.i686)
   - [Darwin](https://github.com/rzab/ostent/releases/download/v0.1.3/Darwin.x86_64)
   - _Expect \*BSD builds surely_

Binaries distributed by [GitHub Releases](https://github.com/rzab/ostent/releases)

Usage
-----

`ostent` accepts optional `-bind` argument to set specific IP and/or port to bind to, otherwise any machine IP and port 8050 by default.

   - `ostent -bind 127.1` # [http://127.0.0.1:8050/](http://127.0.0.1:8050/)
   - `ostent -bind 192.168.1.10:8051` # port 8051
   - `ostent -bind 8052` # any IP, port 8052

`-update` sets collection interval (1 second by default), append `s` for seconds: `0.5s`, `5s` etc.

Run it, it'll give the link(s) to open in a browser.

Running the code
----------------

1. **`git clone https://github.com/rzab/ostent.git`**

2. **`cd ostent`** `# the project directory`

3. **`export GOPATH=$GOPATH:$PWD`** `# the current directory into $GOPATH`

4. **`go get github.com/jteeuwen/go-bindata/go-bindata github.com/rzab/amber`**

5. **`scons`** to generate required `src/ostential/{assets,view}/bindata.devel.go`. These files will contain absolute local paths.
   It's either scons, or run **manually**:
   ```sh
      go-bindata -pkg view   -o src/ostential/view/bindata.devel.go   -tags '!production' -debug -prefix templates.html templates.html
      go-bindata -pkg assets -o src/ostential/assets/bindata.devel.go -tags '!production' -debug -prefix assets         assets/...
   ```

   See [SCons](#scons) on topic.

6. Using [rerun](https://github.com/skelterjohn/rerun), it'll go get the remaining Go dependecies:

	**`go get github.com/skelterjohn/rerun`**

7. **`rerun ostent`**

Go packages
-----------

`[src/]ostential` is the core package.

`[src/]ostent` is the main (_as in [Go Program execution](http://golang.org/ref/spec#Program_execution)_) package:
rerun will find `main.devel.go` file; the other `main.production.go` (used when building with `-tags production`)
is the init code for the distributed binaries: also includes
[goagain](https://github.com/rcrowley/goagain) recovering and self-upgrading via [go-update](https://github.com/inconshreveable/go-update).

`[src/]amberp/amberpp` is templates compiler, used with scons.

SCons
-----

Additional required tools here:
- [Sass](http://sass-lang.com/)
- [react-tools](https://www.npmjs.org/package/react-tools) with [Node.js](http://nodejs.org/)

`scons` makes these **commited to the repo** files:
- `src/ostential/view/bindata.devel.go`
- `src/ostential/assets/bindata.devel.go`
- `assets/css/index.css`
- `assets/js/gen/jscript.js`
- `tmp/jscript.jsx`

If you don't change source files, content re-generated should not differ from the commited.
Whenever amber.templates or assets or style change, you have to re-run `scons`.

`scons build` compiles everything and produces final binary.

The assets
----------

The binaries, to be stand-alone, have the assets (including `templates.html/`) embeded.
Unless you specifically `go build` with `-tags production` (e.g with scons),
they are not embeded for the ease of development:
with `rerun ostent`, asset requests are served from the actual files.
