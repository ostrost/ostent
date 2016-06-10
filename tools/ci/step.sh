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
: "${GIMME_VERSION:=0.2.4}"
: "${GIMME_PATH:=$HOME/bin/gimme}"
: "${GIMME_ENV_PREFIX:=$HOME/.gimme/envs}"
: "${GIMME_VERSION_PREFIX:=$HOME/.gimme/versions}"
export GIMME_ENV_PREFIX GIMME_VERSION_PREFIX

freebsd && DONOTUSE_GIMME=1

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

    if ! eq 1 "$DONOTUSE_GIMME" ; then
        "$GIMME_PATH" "$GO_VERSION" # may be timely
        . "$GIMME_ENV_PREFIX/go$GO_VERSION.env"; go env >&2 #source here & verbose to &2
    fi

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
    # unconditionally install gimme(1); travis env most definitely does not have it
    # if ! eq 1 "$DONOTUSE_GIMME" ; then # travis always uses gimme
    mkdir -p "$(dirname "$GIMME_PATH")"
    curl -sSL -o "$GIMME_PATH" https://github.com/travis-ci/gimme/raw/v"$GIMME_VERSION"/gimme # timely
    chmod +x "$GIMME_PATH"
}

install_2() {
    local REPOSLUG="$1"

    # if ! eq 1 "$DONOTUSE_GIMME" ; then # travis always uses gimme
    "$GIMME_PATH" "$GO_VERSION" # timely
    . "$GIMME_ENV_PREFIX/go$GO_VERSION.env"; go env >&2 #source here & verbose to &2

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
    # It's a runner target so the env/state is clean but prepped with before_script.
    glide install # partial Gmake init
    Gmake test
}
cideploy() { # Gmake deploy
    # It's a runner target so the env/state is clean but prepped with before_script.
    glide install # partial Gmake init

    # For a runner, bootstrapping must have been done
    # For Travis CI, before_deploy_{1,2} bootstrap the 32-bit cross building

    before_deploy_3
    before_deploy_4

    local tag
    tag=$(git describe --tags --abbrev=0) # literal tag, should be in "v..." form

    local release=/home/release/bin/github-release
    # local tagsansv=${tag##v}
    # "$release" release \
    #            --tag "$tag" \
    #            --name "ostent $tagsansv" \
    #            --description ' ' \
    #            --draft \
    #            --pre-release

    for filename in "$DPL_DIR"/* ; do
        "$release" upload \
                   --tag "$tag" \
                   --file "$filename" \
                   --name "TESTING-$(basename "$filename")"
        # The "TESTING-" prefix until we done testing this.
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
        # if ! eq 1 "$DONOTUSE_GIMME" ; then # travis always uses gimme
        "$GIMME_PATH" "$GO_BOOTSTRAPVER" >/dev/null # timely
    fi
}

before_deploy_2() {
    if ! darwin ; then
        # if ! eq 1 "$DONOTUSE_GIMME" ; then # travis always uses gimme
        GOROOT_BOOTSTRAP="$(ls -d "$GIMME_VERSION_PREFIX/go$GO_BOOTSTRAPVER".*.amd64)" \
        Gmake boot32
    fi
}

before_deploy_3() {
    if ! darwin ; then
        Gmake all32
    fi
}

before_deploy_4() {
    DPL_DIR=$(eval echo "$DPL_DIR")
    mkdir -p "$DPL_DIR"

    local u; u=$(uname)
    before_deploy_fptar "$u"
    before_deploy_fptar "$u" 32

    (cd "$DPL_DIR" && shasum -a 256 ./*.tar.xz >CHECKSUM."$u".SHA256)
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

"$@" # The last line to dispatch. $1 is ought to be a func name.
