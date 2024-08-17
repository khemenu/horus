#!/usr/bin/env sh
set -e

case "$1" in
    horus|hr)
        exec "$@"
        ;;
    *)
        exec horus "$@"
        ;;
esac
