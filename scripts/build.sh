#!/bin/bash -e

set -euf -o pipefail

rm -rf vendor/github.com/ankyra/escape-core
cp -r deps/_/escape-core/ vendor/github.com/ankyra/escape-core
rm -rf vendor/github.com/ankyra/escape-core/vendor/

user_id=$(id -u $(whoami))

docker run --rm \
    -v "$PWD":/go/src/github.com/ankyra/escape-registry \
    -w /go/src/github.com/ankyra/escape-registry \
    golang:1.9.0 bash -c "(useradd --uid $user_id builder || true) && su builder -p -c \"/usr/local/go/bin/go build -v\""
