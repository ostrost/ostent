#!/bin/sh
test -n "$BASH" -o -n "$KSH_VERSION" -o -n "$ZSH_VERSION" &&
set -o pipefail 2>/dev/null
set -eu
if test x${TRAVIS:-false} == xtrue ; then
    set +u # non-strict unset variables use in CI build script
fi

GO_BOOTSTRAPVER=go1.4.3
: ${DPL_DIR:=$(git rev-parse --show-toplevel)/deploy}

eq()      { test x$1 == x$2 ;}
linux()   { eq "$1" Linux   ;}
darwin()  { eq "$1" Darwin  ;}
freebsd() { eq "$1" FreeBSD ;}

: ${MAKE:=make}
Gmake() {
    if eq $MAKE make && hash gmake 2>/dev/null ; then
        MAKE=gmake
    fi
    $MAKE "$@"
}

# Following functions of this script is expected to be executed sequentially.
# The split is so that each function must end with one timely action.

install_1() {
    if hash gvm 2>/dev/null ; then
        gvm get
    else
        curl -sSL https://github.com/moovweb/gvm/raw/master/binscripts/gvm-installer |
        bash
    fi
}

install_2() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"

    if darwin && ! eq $GOVER tip ; then
        GO_BINARY_PATH=~/.gvm/archive///////$GOVER.darwin-amd64-osx10.8.tar.gz
        GO_BINARY_URL=https://golang.org/dl/$GOVER.darwin-amd64.tar.gz
        test -f $GO_BINARY_PATH ||
        curl --silent --show-error --fail --location --output $GO_BINARY_PATH $GO_BINARY_URL
    fi
}

install_3() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"

    source ~/.gvm/scripts/gvm
    gvm version
    gvm install $GOVER --binary # || gvm install $GOVER
}

# Nothing timely here, but it's the last install step.
install_4() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"
    local REPOSLUG="$2" # The "owner/repo" form.

    gvm use $GOVER
    gvm list

    mkdir -p ~/gopath/src
    mv ~/build ~/gopath/src/github.com # ~/build is cwd
    cd ~/gopath/src/github.com/$REPOSLUG # new cwd

    export GOPATH=~/gopath:$GOPATH # NB
    export PATH=~/gopath/bin:$PATH

    go version
    go env
}

before_deploy_1() {
    if ! darwin ; then
        gvm install $GO_BOOTSTRAPVER --binary || true
        (
            gvm use $GO_BOOTSTRAPVER
            gvm pkgset list
        )
    fi
}

before_deploy_2() {
    if ! darwin ; then
        Gmake boot32 GOROOT_BOOTSTRAP=~/.gvm/gos/$GO_BOOTSTRAPVER
    fi
}

before_deploy_3() {
    local uname=$(uname)
    local arch=$(uname -m)

    if ! darwin $uname ; then
        Gmake all32
    fi
}

before_deploy_4() {
    local uname=${1:-$(uname)}

    mkdir -p "$DPL_DIR"
    before_deploy_fptar $uname
    before_deploy_fptar $uname 32

    local shacommand=shasum\ -a\ 256
    if ! hash shasum 2>/dev/null ; then
        shacommand=sha256\ -r # freebsd
    fi
    (
        cd "$DPL_DIR" || exit 1
        find . -type f \! -name CHECKSUM.\* | sed 's,^\./,,' |
        xargs $shacommand >CHECKSUM."$uname".SHA256
    )
}

before_deploy_fptar() {
    local uname=${1:-$(uname)}
    local arch=${2:-$(uname -m)}

    if freebsd $uname ; then
        if eq $arch x86_64 ; then
            arch=amd64
        fi
    elif eq $arch == amd64 ; then
        arch=x86_64
    fi

    local binary=~/gopath/bin/ostent
    if eq $arch 32 ; then
        binary=~/gopath/bin/ostent.32
        arch=i686
        if freebsd $uname ; then
            arch=i386
        fi
    fi

    if darwin $uname && ! eq $arch x86_64 ; then
        return # Only 64-bit builds for darwin
    fi

    local prefix=/usr
    if ! linux $uname; then
        prefix=/usr/local
    fi

    local tarball="$DPL_DIR"/$uname-$arch.tar.xz
    if test -e "$tarball" ; then
        echo File already exists: "$tarball" >&2
        exit 1
    fi
    local tmpsubdir=$(mktemp -d tmpstage.XXXXXXXX) || exit 1 # in cwd
    trap 'rm -rf "$PWD"/'"$tmpsubdir" EXIT
    (
        cd "$tmpsubdir" || exit 1

        # umask 022 # MIND UMASK
        install -m 755 -d . ./$prefix/bin
        install -m 755 -p $binary ./$prefix/bin/ostent
        find . -type d | xargs touch -r $binary

        local ownerargs=--owner=0\ --group=0
        if freebsd $uname ; then
            ownerargs=--uid=0\ --gid=0
        elif darwin $uname ; then
            ownerargs= # No way to specify owners in darwin
        fi

        echo Packing $uname-$arch >&2
        tar Jcf "$tarball" --numeric-owner $ownerargs .
    )
    rm -rf "$tmpsubdir"
    # trap EXIT # clear the trap
}

prior_to_deploy() {
    set +e # off fatal errors for travis-dpl
}

"$@" # The last line to dispatch. $1 is ought to be a func name.
