#!/bin/sh
# strict mode
test -n "$BASH" -o -n "$KSH_VERSION" -o -n "$ZSH_VERSION" &&
    set -o pipefail 2>/dev/null
set -eu
set -x

EXITunless_type() {
	  type "$1" >/dev/null || exit 0
}

SCRIPT_webpack() {
    EXITunless_type webpack
    webpack "$input" "$output"
}

SCRIPT_webpack_ugly() {
    EXITunless_type uglifyjs
    SCRIPT_webpack
	  uglifyjs --compress --mangle --output "$output" "$output"
}

# first argument must be the name of a function: SCRIPT_...
"$@"
