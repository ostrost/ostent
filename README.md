[![Build status](https://secure.travis-ci.org/ostrost/ostent.png?branch=master)](https://travis-ci.org/ostrost/ostent)
[![Sourcegraph](https://sourcegraph.com/api/repos/github.com/ostrost/ostent/.badges/status.png)](https://sourcegraph.com/github.com/ostrost/ostent)
[![GoDoc](https://godoc.org/github.com/ostrost/ostent?status.svg)](https://godoc.org/github.com/ostrost/ostent)

`ostent` collects and displays system metrics and optionally relays to Graphite and/or InfluxDB.

The displaying part ([**demo**](http://demo.ostrost.com/)) is interactive and customizable.

![Screenshot](https://www.ostrost.com/ostent/screenshot.png)

The metrics:
- **Collected and exported**:
  - RAM, swap usage
  - Interfaces bytes, packets, errors ins and outs
  - CPU usage
  - Disk usage
  - Load average
- **On display only**:
  - System's OS, IP and uptime
  - Processes top
  - vagrant global-status

The exporting to Graphite and InfluxDB is kept on par with [collectd](https://collectd.org/)
[reporting](https://collectd.org/wiki/index.php/Plugin:Write_Graphite) to Graphite with `StoreRates true`,
although the metrics naming is slightly different.

Running
-------

ostent a single executable without dependecies, no extra files required (everything is builtin).
Drop it in and just run; being root is unnecesary. There're [flags](#usage) if you have to.

[Run the code](#running-the-code) if you want to, otherwise grab a binary.

Install Release binaries
========================

These binaries self-upgrade whenever there's new stable release and
distributed by [GitHub Releases](https://github.com/ostrost/ostent/releases).

Install & run with `curl -sSL https://github.com/ostrost/ostent/raw/master/ostent.sh | sh`

Platforms

   - Linux [64-bit](https://github.com/ostrost/ostent/releases/download/v0.2.0/Linux.x86_64) | [32-bit](https://github.com/ostrost/ostent/releases/download/v0.2.0/Linux.i686)
   - [Mac OS X](https://github.com/ostrost/ostent/releases/download/v0.2.0/Darwin.x86_64) (64-bit)

FreeBSD (10 amd64 and i386 probably) is to be published with a new release.
The master code is runnable already.

Usage
-----

```
Usage of ostent:
  -bind=:8050: Bind address
  -update=1s: Collection interval

  -sendto-graphite=: Graphite server address
  -graphite-refresh=10s: Graphite refresh interval

  -sendto-influxdb=: InfluxDB server URL
  -influxdb-database="ostent": InfluxDB database
  -influxdb-password="": InfluxDB password
  -influxdb-refresh=10s: InfluxDB refresh interval
  -influxdb-username="": InfluxDB username
```

Unless `-bind` (`-b` for short) is set, ostent binds to `*:8050`.
The bind and Graphite addresses are specified like `IP[:port]`
(default ports being 8050 and 2003 respectively).
InfluxDB server must be specified as an URL `http://ADDRESS`.
An interval is a number and a unit: `s` for seconds, `m` for minutes etc.

Here's how it goes:

```
$ ostent                                     ________________
[ostent]    -------------                   < Spot the links >
[ostent]  / server ostent \                  ----------------
[ostent] +------------------------------+           \   ^__^
[ostent] | http://127.0.0.1:8050        |            \  (oo)\_______
[ostent] |------------------------------|               (__)\       )\/\
[ostent] | http://192.168.1.2:8050      |                   ||----w |
[ostent] +------------------------------+                   ||     ||
```

Running the code
----------------

Have your GOPATH environment set,
[gvm](https://github.com/moovweb/gvm) is a must.

1. **`go get github.com/ostrost/ostent`**
2. `ostent` to run.

For rebuilding the code and assets:

1. Find `src/github.com/ostrost/ostent` directory in GOPATH.
2. Run `make init` once and later for packages update (think `go get -u`)
3. `make` or `make al` when `rerun` does rebuilding.

Repeat 3. every time sources (esp. assets) change.
[rerun](https://github.com/skelterjohn/rerun) does live-reloading run:
`rerun github.com/ostrost/ostent`

**For a fork**, to preserve import paths and packages namespace,
clone your fork as if it was `github.com/ostrost/ostent` package for Go:

1. `go get github.com/ostrost/ostent`
2. Find `src/github.com/ostrost/ostent` directory in GOPATH.
3. Replace it with you fork clone.
4. Continue with rebuilding steps above.

Make
----

`make` rebuilds these **commited to the repo** files:
- `share/templates/bindata.*.go`
- `share/assets/bindata.*.go`
- `share/assets/js/devel/milk/*.js`
- `share/assets/js/devel/gen/*.js`
- `share/templates/*.html`
- `share/assets/css/*.css`
- `share/tmp/*.jsx`

If you don't change source files, content re-generated
should not differ from the commited. Whenever
share/{amber.templates,style,coffee} modified,
you have to re-make.

Additional tools required for assets rebuilding:
- [Sass](http://sass-lang.com/install)
- [react-tools](https://www.npmjs.org/package/react-tools)
- [coffee-script](https://www.npmjs.com/package/coffee-script)
- [requirejs](https://www.npmjs.org/package/requirejs)

The main package
----------------

`github.com/ostrost/ostent` has two main.go files:
rerun will find `main.devel.go`; the other `main.production.go`
(used when building with `-tags production`) is the init code for
the distributed binaries: also includes
[goagain](https://github.com/rcrowley/goagain) recovering and
self-upgrading via [go-update](https://github.com/inconshreveable/go-update).

The assets
----------

The binaries, to be stand-alone, have the assets and templates embeded.
Unless you specifically build with `-tags production` (e.g with make),
the content is not embeded for the ease of development:
with `rerun`, asset requests are served from the actual files.
Production-built `ostent restore-assets` can be used to copy (extract) assets on disk.
