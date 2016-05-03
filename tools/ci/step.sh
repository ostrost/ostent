#!/bin/sh
test -n "${BASH:-}" -o -n "${KSH_VERSION:-}" -o -n "${ZSH_VERSION:-}" &&
set -o pipefail 2>/dev/null
set -eu
eq() { test "x$1" = "x$2"; }
if eq "${TRAVIS:-}" true ; then
    set +u # non-strict unset variables use in CI build script
fi

if ! eq "${TRAVIS:-}" true ; then
    set -x #while debugging
fi

GO_BOOTSTRAPVER=go1.5.4
: "${GIT_TOPLEVEL:=$(git rev-parse --show-toplevel)}"
: "${DPL_DIR:=$GIT_TOPLEVEL/deploy}"

linux()   { eq "${1:-$(uname)}" Linux   ;}
darwin()  { eq "${1:-$(uname)}" Darwin  ;}
freebsd() { eq "${1:-$(uname)}" FreeBSD ;}

: "${MAKE:=make}"
Gmake() {
    if eq "$MAKE" make && hash gmake 2>/dev/null ; then
        MAKE=gmake
    fi
    "$MAKE" "$@"
}

# Following functions of this script is expected to be executed sequentially.
# The split is so that each function must end with one timely action.

: "${GO_VERSION:=1.6.2}"
: "${GIMME_PATH:=~/bin/gimme}"
: "${GIMME_ENV_PREFIX:=~/.gimme/envs}"
: "${GIMME_VERSION_PREFIX:=~/.gimme/versions}"
export GIMME_ENV_PREFIX GIMME_VERSION_PREFIX

# before_script is executed by gitlab-runner
before_script() {
    # required in environ: $GO_VERSION

    local d ownername reponame
    # d=/home/gitlab-runner\
    # /builds/${runner_id:-deadbeef}/${CI_PROJECT_ID:-0}/$ownername/$reponame
    d=${CI_PROJECT_DIR:-${TRAVIS_REPO_SLUG:-$GIT_TOPLEVEL}}
    ownername=$(basename "$(dirname "$d")")
    reponame=$(basename "$d")
    export GOPATH="$HOME/gopath-$ownername-$reponame"

    "$GIMME_PATH" "$GO_VERSION" # may be timely
    . "$GIMME_ENV_PREFIX/go$GO_VERSION.env"; go env >&2 #source here & verbose to &2
    PATH=''/home/glide/bin:"$PATH"; export PATH
    glide --version >&2

    package=$(glide name) # $(awk '/^package: / { print $2 }' glide.yaml)
    # if eq "$package" ''; then
    #     # || git remote show -n origin|x
    #     # || git config remote.origin.url
    #     local h
    #     h=$(git ls-remote --get-url)
    #     h=${h#*://} # remove prefix
    #     h=${h#*@}   # remove prefix
    #     h=${h%:*}   # remove suffix
    #     h=${h%%/*}  # remove suffix, greedily
    #     package="$h/$ownername/$reponame"
    # fi
    export package # [exported] $package used with Gmake

    local symlink readlink destlink
    symlink="$GOPATH/src/$package"
    readlink=$(readlink "$symlink") || true
    destlink="$PWD/"
    if ! eq "$readlink" "$destlink" ; then
        { #debug
            ls -ld "$symlink" || true
            echo link read: "$readlink"
            echo link should: "$destlink"
        }
        rm -rf "$symlink"
        mkdir -p "$(dirname "$symlink")"
        ln -s "$destlink" "$symlink"
    fi

    local import
    import=$(cd "$symlink" && go list -f '{{.ImportPath}}')
    if ! eq "$import" "$package" ; then
        echo "Import path skewed: package=$package go:$import" >&2
        exit 1
    fi
}

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

    if darwin && ! eq "$GOVER" tip ; then
        local GO_BINARY_PATH GO_BINARY_URL
        GO_BINARY_PATH=~/.gvm/archive///////"$GOVER".darwin-amd64-osx10.8.tar.gz
        GO_BINARY_URL=https://golang.org/dl/"$GOVER".darwin-amd64.tar.gz
        test -f "$GO_BINARY_PATH" ||
        curl --silent --show-error --fail --location --output "$GO_BINARY_PATH" "$GO_BINARY_URL"
    fi
}

install_3() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"

    . ~/.gvm/scripts/gvm #source here
    gvm version
    gvm install "$GOVER" --binary # || gvm install "$GOVER"
}

# Nothing timely here, but it's the last install step.
install_4() {
    local GOVER="$1" # go version in form of "goX.Y[.Z]"
    local REPOSLUG="$2" # The "owner/repo" form.

    gvm use "$GOVER"
    gvm list

    mkdir -p ~/gopath/src
    mv ~/build ~/gopath/src/github.com # ~/build is cwd
    cd ~/gopath/src/github.com/"$REPOSLUG" # new cwd

    export GOPATH=~/gopath # NB
    export PATH=~/gopath/bin:$PATH

    go version
    go env
}

cibuild() {
    Gmake init
    Gmake --always-make all
}
citest() {
    # citest is a target so the env/state is clean but prepped with before_script.
    glide install # partial Gmake init
    Gmake test
}
cideploy() { # Gmake deploy
    # cideploy is a target so the env/state is clean but prepped with before_script.
    glide install # partial Gmake init

    # before_deploy_1
    # before_deploy_2
    # before_deploy_3
    if ! darwin ; then
        # bootstrapping must have been done
        Gmake all32
    fi
    before_deploy_4

    local tag=$(git describe --tags --abbrev=0) # literal tag, should be in "v..." form
    local tagsansv=${tag##v}

    local release=/home/release/bin/github-release
    # "$release" release \
    #            --tag "$tag" \
    #            --name "ostent $tagsansv" \
    #            --description ' ' \
    #            --draft \
    #            --pre-release

    for filename in "$DPL_DIR"/* ; do
        "$release" upload \
                   --tag "$tag" \
                   --name $(basename "$filename") \
                   --file test-"$filename" # NB
    done
}

maketest() {
    local ps=${testpackages:-./...}
    if (eq "$ps" '' || eq "$ps" ./...) && hash glide 2>/dev/null ; then
        ps=$(glide novendor | grep -v '^\./builds/')
    fi

    echo "$ps" | xargs go vet

    local import="${GOPATH%%:*}/src/$package" # go list -f {{.Dir}} "$package"
    (cd "$import" && echo "$ps" | xargs go test -v)
}

covertest() {
    local sp=${testpackage:-./...}
    if eq "$sp" ./... ; then
        sp=${package:-$sp}
    fi
    go test -coverprofile=coverage.out -covermode=count -v "$sp"
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
    if ! darwin ; then
        Gmake all32
    fi
}

before_deploy_4() {
    local uname=${1:-$(uname)}

    mkdir -p "$DPL_DIR"
    before_deploy_fptar "$uname"
    before_deploy_fptar "$uname" 32

    (cd "$DPL_DIR" && shasum -a 256 ./*.tar.xz >CHECKSUM."$uname".SHA256)
}

before_deploy_fptar() {
    local uname=${1:-$(uname)}
    local arch=${2:-$(uname -m)}

    if freebsd "$uname" ; then
        if eq "$arch" x86_64 ; then
            arch=amd64
        fi
    elif eq "$arch" amd64 ; then
        arch=x86_64
    fi

    local binary="$GOPATH"/bin/ostent
    if eq "$arch" 32 ; then
        binary="$GOPATH"/bin/ostent.32
        arch=i686
        if freebsd "$uname" ; then
            arch=i386
        fi
    fi

    if darwin "$uname" && ! eq "$arch" x86_64 ; then
        return # Only 64-bit builds for darwin
    fi

    local prefix=/usr
    if ! linux "$uname"; then
        prefix=/usr/local
    fi

    local tarball="$DPL_DIR/$uname-$arch".tar.xz
    if test -e "$tarball" ; then
        echo File already exists: "$tarball" >&2
        exit 1
    fi
    local tmpsubdir
    tmpsubdir=$(mktemp -d tmpstage.XXXXXXXX) || exit 1 # in cwd
    trap 'rm -rf "$PWD"/'"$tmpsubdir" EXIT
    (
        cd "$tmpsubdir" || exit 1

        # umask 022 # MIND UMASK
        install -m 755 -d . ./$prefix/bin
        install -m 755 -p $binary ./$prefix/bin/ostent
        find . -type d -print0 | xargs -0 touch -r $binary

        local ownerargs=--owner=0\ --group=0
        if freebsd "$uname" ; then
            ownerargs=--uid=0\ --gid=0
        elif darwin "$uname" ; then
            ownerargs='' # No way to specify owners in darwin
        fi

        echo "Packing $uname-$arch" >&2
        tar Jcf "$tarball" --numeric-owner $ownerargs .
    )
    rm -rf "$tmpsubdir"
    # trap '' EXIT # clear the trap
}

prior_to_deploy() {
    set +e # off fatal errors for travis-dpl
}

"$@" # The last line to dispatch. $1 is ought to be a func name.
