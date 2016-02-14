#!/bin/sh
test -n "$BASH" -o -n "$KSH_VERSION" -o -n "$ZSH_VERSION" &&
set -o pipefail 2>/dev/null
set -e # not -u
set +u # non-strict unset variables use in CI build script

set -x # debug

GO_BOOTSTRAPVER=1.4.3
UNAME=$(uname)
ARCH=$(uname -m)
: ${MAKE:=make}

Gmake() {
    if test x$MAKE == xmake && hash gmake 2>/dev/null ; then
        MAKE=gmake
    fi
    $MAKE "$@"
}

# Following functions of this script is expected to be executed sequentially.
# The split is so that each function must end with one timely action.

install1st() {
    if hash gvm 2>/dev/null ; then
        gvm get
    else
        curl -sSL https://github.com/moovweb/gvm/raw/master/binscripts/gvm-installer |
        bash
    fi
}

install2nd() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"
    local OSXOS="$2" # "osx" if host is a mac

    if test x$OSXOS == xosx -a x$GOVER != xtip ; then
        GO_BINARY_PATH=~/.gvm/archive///////$GOVER.darwin-amd64-osx10.8.tar.gz
        GO_BINARY_URL=https://golang.org/dl/$GOVER.darwin-amd64.tar.gz
        test -f $GO_BINARY_PATH ||
        curl --silent --show-error --fail --location --output $GO_BINARY_PATH $GO_BINARY_URL
    fi
}

install3rd() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"

    source ~/.gvm/scripts/gvm
    gvm version
    gvm install $GOVER --binary # || gvm install $GOVER
}

# Nothing timely here, but it's the last install step.
install4th() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"
    local REPOSLUG="$2" # The "owner/repo" form.

    gvm use $GOVER
    gvm list

    cd
    mkdir -p gopath/src
    mv build gopath/src/github.com
    cd gopath/src/github.com/$REPOSLUG

    export GOPATH=~/gopath:$GOPATH # NB
    export PATH=~/gopath/bin:$PATH
    export CC=clang CXX=clang++

    go version
    go env
}

before-deploy1st() {
    local OSXOS="$1" # "osx" if host is a mac

    mkdir -p deploy/ # NB

    if test x$OSXOS != xosx ; then
        gvm install $GO_BOOTSTRAPVER --binary || true
        (gvm use $GO_BOOTSTRAPVER; gvm pkgset list)
    fi
}

before-deploy2nd() {
    local OSXOS="$1" # "osx" if host is a mac

    if test x$OSXOS != xosx ; then
        Gmake boot32 GOROOT_BOOTSTRAP=~/.gvm/gos/$GO_BOOTSTRAPVER
    fi
}

before-deploy3rd() {
    local OSXOS="$1" # "osx" if host is a mac

    if test x$OSXOS != xosx ; then
        Gmake all32
        cp -p $HOME/gopath/bin/ostent.32 deploy/$UNAME.i686
    fi
    cp -p $HOME/gopath/bin/ostent deploy/$UNAME.$ARCH

    set +e # NB off fatal errors for travis-dpl
}

"$@" # The last line to dispatch. $1 is ought to be a func name.
