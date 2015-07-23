#!/bin/sh
version="$1"
if test x"$version" == x; then
    echo Usage: $0 version >&2
    exit 64
fi

subdir=bootstrap-"$version"
archive=bootstrap-"$version".zip
test -e "$archive" || curl -Lo "$archive" https://github.com/twbs/bootstrap/archive/v"$version".zip

rm -rf "$subdir"
unzip "$archive" "$subdir"/less/\*
