#!/bin/sh

CURRENT=$(cd "$(dirname "$0")" && pwd)
docker run --rm -it \
    -v "$CURRENT":/go/src/github.com/shogo82148/go-prove \
    -w /go/src/github.com/shogo82148/go-prove golang:1.9 "$@"
