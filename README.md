# Ostent [![Travis CI][2]][1]
[1]: https://travis-ci.org/ostrost/ostent
[2]: https://travis-ci.org/ostrost/ostent.svg?branch=master

Ostent collects metrics to display and report to InfluxDB, Graphite, Librato.

The interactive display UI ([**demo**](https://demo.ostrost.com/)):

![Screenshot](https://www.ostrost.com/ostent/screenshot.png)

System metrics collected and reported:
- RAM, swap usage
- CPU usage, load average
- Disk space usage in bytes and inodes
- Network ins and outs in bytes, packets, drops and errors

The processes top is on-display only.

## Install

Ostent is a single executable.
[Release tarball](https://github.com/ostrost/ostent/releases)
has the binary &mdash; download and extract in one go:

```sh
curl -L https://github.com/ostrost/ostent/releases/download/v0.6.2/`uname`-`uname -m`.tar.xz | tar Jxf -
```

This will place executable in `./usr/**/bin/ostent`.
For system-wide install use `sudo tar Jxf - -C / <<<...`.

Platforms:

   - Linux
   - FreeBSD
   - Mac OS X

## Usage

```
$ ostent -h
Ostent is a server tool to collect, display and report system metrics.

Usage:
  ostent [flags]

Flags:
      --bind string         server bind address (default "")
      --bind-port int       server bind port (default 8050)
      --config string       config filename (default "$HOME/.ostent.toml")
      --interval duration   metrics collection interval (default 10s)
      --log-requests        log server requests (default false)
      --upgrade-checks      periodic upgrade checks (default true)
      --version             print version and exit
```

## Running the code

1. `go get github.com/ostrost/ostent`
2. `ostent` to run.

See also [Two kind of builds](#two-kinds-of-builds).

### Rebuilding

1. `cd $(go list -f {{.Dir}} github.com/ostrost/ostent)`
2. `make init` once.
3. `npm install` once, optional, sets up assets and template rebuilding.
4. `gulp watch` or `make` after changes.

`make` rebuilds these **commited to the repo** files:
- `share/assets/bindata.*.go`
- `share/assets/css/*.css`
- `share/assets/js/*/*.js`
- `share/templates/*.html`
- `share/templates/bindata.*.go`
- `share/js/*.jsx`

If you don't change source files, content re-generated
should not differ from the commited.

`gulp watch`

- watches share/{js,style,templatesorigin} and rebuilds dependants on changes
- does live-reloading `ostent` code run
- acceps all ostent flags e.g. `gulp watch -b 127.0.0.1:8080`

### Two kinds of builds

Standalone and release binaries produced by `make` (or `go get -tags bin`)
include embeded template and assets.

Non-bin builds made by `gulp watch` and `go get`
- serve assets and use template from actual files
- have a set of flags facilitating debugging etc.
