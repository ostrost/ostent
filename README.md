# ostent [![GoDoc][2]][1] [![Travis CI][4]][3]
[1]: https://godoc.org/github.com/ostrost/ostent
[2]: https://godoc.org/github.com/ostrost/ostent?status.svg
[3]: https://travis-ci.org/ostrost/ostent
[4]: https://travis-ci.org/ostrost/ostent.svg?branch=master

Ostent collects system metrics to display and relay to

- Graphite
- InfluxDB
- Librato

The interactive UI ([**demo**](https://demo.ostrost.com/)):

![Screenshot](https://www.ostrost.com/ostent/screenshot.png)

The metrics collected and exported:
- RAM, swap usage
- CPU usage, load average
- Disk space usage in bytes and inodes
- Network ins and outs in bytes, packets, drops and errors

The processes top is on-display only.

The exporting is kept on par with [collectd](https://collectd.org/)
[reporting](https://collectd.org/wiki/index.php/Plugin:Write_Graphite),
although the metrics naming is slightly different.

## Install

Ostent is a single executable.
[Release tarball](https://github.com/ostrost/ostent/releases)
has the binary &mdash; download and extract in one go:

```sh
tar Jxf - <<<$(curl -L https://github.com/ostrost/ostent/releases/download/v0.6.1/`uname`-`uname -m`.tar.xz)
```

This will place executable in `./usr/**/bin/ostent`.
For system-wide install use `sudo tar Jxf - -C / <<<...`.

Platforms:

   - Linux
   - FreeBSD
   - Mac OS X

## Usage

```
To continuously export collected metrics, use --graphite, --influxdb and/or --librato.
Use multiple flags and/or use comma separated endpoints for the same kind. E.g.:
      --graphite 10.0.0.1,10.0.0.2:2004\?delay=30s
      --influxdb http://10.0.0.3\?delay=60s
      --librato \?email=EMAIL\&token=TOKEN
ostent -h lists supported parameters and defaults.

Flags:
  -b, --bind address            Bind address (default :8050)
      --graphite endpoints      Graphite exporting endpoints
      --influxdb endpoints      InfluxDB exporting endpoints
      --librato parameters      Librato exporting parameters
      --max-delay delay         Maximum for display delay (default 10m)
  -d, --min-delay delay         Collection and display minimum delay (default 1s)
      --version                 Print version and exit
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
