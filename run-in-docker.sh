#!/bin/sh

CURRENT=$(cd "$(dirname "$0")" && pwd)
docker run --rm -it \
    -e GO111MODULE=on \
    -e "GOOS=${GOOS:-linux}" -e "GOARCH=${GOARCH:-amd64}" -e "CGO_ENABLED=${CGO_ENABLED:-0}" \
    -v "$CURRENT":/go/src/github.com/shogo82148/go-prove \
    -v "$CURRENT/.mod":/go/pkg/mod \
    -w /go/src/github.com/shogo82148/go-prove golang:1.13.6 "$@"
