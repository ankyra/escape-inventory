#!/bin/bash -e

set -e -o pipefail


PLATFORMS="darwin linux"
ARCHS="386 amd64"

echo "$INPUT_credentials" > service_account.json
gcloud auth activate-service-account --key-file service_account.json

for GOOS in $PLATFORMS; do
    for ARCH in $ARCHS; do
        target="escape-inventory-v$INPUT_inventory_version-$GOOS-$ARCH.tgz"
        echo "Building $target"
        if [ ! -f $target ] ; then
            docker run --rm -v "$PWD":/go/src/github.com/ankyra/escape-inventory \
                            -w /go/src/github.com/ankyra/escape-inventory \
                            -e GOOS=$GOOS \
                            -e GOARCH=$ARCH \
                            golang:1.8 go build -v -o escape-inventory-$GOOS-$ARCH
            mv escape-inventory-${GOOS}-${ARCH} escape-inventory
            tar -cvzf ${target} escape-inventory
            rm escape-inventory
        else
            echo "File $target already exists"
        fi
        gcs_target="gs://$INPUT_bucket/escape-inventory/$INPUT_inventory_version/$target"
        echo "Copying to $gcs_target"
        gsutil cp "$target" "$gcs_target"
        echo "Setting ACL on $gcs_target"
        gsutil acl ch -u AllUsers:R "$gcs_target"
    done
done
