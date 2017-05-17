#!/bin/bash -e

set -euf -o pipefail

user_id=$(id -u $(whoami))

docker run --rm \
    -v "$PWD":/go/src/github.com/ankyra/escape-registry \
    -w /go/src/github.com/ankyra/escape-registry \
    golang:1.8 bash -c "(useradd --uid $user_id builder || true) && su builder -p -c \"/usr/local/go/bin/go test -cover -v \\\$(/usr/local/go/bin/go list ./... | grep -v -E 'vendor' )\""

