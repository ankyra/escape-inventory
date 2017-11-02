#!/bin/bash -e

set -euf -o pipefail

rm -rf vendor/github.com/ankyra/escape-core
cp -r deps/_/escape-core/ vendor/github.com/ankyra/escape-core
rm -rf vendor/github.com/ankyra/escape-core/vendor/

docker rm src || true
docker create -v /go/src/github.com/ankyra/ --name src golang:1.9.0 /bin/true
docker cp "$PWD" src:/go/src/github.com/ankyra/tmp
docker run --rm --volumes-from src \
    -w /go/src/github.com/ankyra/ \
    golang:1.9.0 mv tmp escape-inventory
docker run --rm \
    --volumes-from src \
    -w /go/src/github.com/ankyra/escape-inventory \
    golang:1.9.0 bash -c "go build"
docker cp src:/go/src/github.com/ankyra/escape-inventory/escape-inventory escape-inventory
docker rm src
