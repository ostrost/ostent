#!/bin/sh -e
set -e # yeah, won't ignore errors

DEST="${DEST:-$HOME/bin/ostent}" # change if you wish. the directory must be writable for ostent to self-update

if ! test -e "$DEST" ; then
    VERSION=v0.1.3
    URL="https://OSTROST.COM/ostent/releases/latest/$(uname -sm)/ostent"
    URL="https://github.com/rzab/ostent/releases/download/$VERSION/$(uname -sm | tr \  .)"

    curl -sSL --create-dirs -o "$DEST" "$URL"
    chmod +x "$DEST"
fi

for arg in in "$@" ; do
    test "x$arg" == x-norun &&
    exit # Ok, just install, no run
done

exec "$DEST" -updatelater "$@"
