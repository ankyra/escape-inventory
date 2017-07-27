#!/bin/bash -e

set -e -o pipefail


PLATFORMS="darwin linux"
ARCHS="386 amd64"

echo "$INPUT_credentials" > service_account.json
gcloud auth activate-service-account --key-file service_account.json

for GOOS in $PLATFORMS; do
    for ARCH in $ARCHS; do
        target="escape-registry-v$INPUT_version-$GOOS-$ARCH.tgz"
        echo "Building $target"
        if [ ! -f $target ] ; then
            docker run --rm -v "$PWD":/go/src/github.com/ankyra/escape-registry \
                            -w /go/src/github.com/ankyra/escape-registry \
                            -e GOOS=$GOOS \
                            -e GOARCH=$ARCH \
                            golang:1.8 go build -v -o escape-registry-$GOOS-$ARCH
            mv escape-registry-${GOOS}-${ARCH} escape-registry
            tar -cvzf ${target} escape-registry
            rm escape-registry
        else
            echo "File $target already exists"
        fi
        gcs_target="gs://$INPUT_bucket/escape-registry/$INPUT_version/$target"
        gsutil cp "$target" "$gcs_target"
        gsutil acl ch -u AllUsers:R "$gcs_target"
    done
done
