#!/bin/bash -e

set -e -o pipefail


PLATFORMS="linux"
ARCHS="amd64"

BASE_DIR=$(dirname "$(readlink -f "$0")")
SRC_DIR=$(readlink -f "${BASE_DIR}/../")

echo "$INPUT_credentials" > service_account.json

if [ "$INPUT_do_upload" = "1" ] ; then 
    gcloud auth activate-service-account --key-file service_account.json
fi

for GOOS in $PLATFORMS; do
    for ARCH in $ARCHS; do
        filename="escape-inventory-v$INPUT_inventory_version-$GOOS-$ARCH.tgz"
        target="${SRC_DIR}/${filename}"
        echo "Building $target"
        if [ ! -f $target ] ; then
            docker run --rm -v "$SRC_DIR":/go/src/github.com/ankyra/escape-inventory \
                            -w /go/src/github.com/ankyra/escape-inventory \
                            -e GOOS=$GOOS \
                            -e GOARCH=$ARCH \
                            golang:1.8 go build -v -o escape-inventory-$GOOS-$ARCH
            mv escape-inventory-${GOOS}-${ARCH} "${SRC_DIR}/escape-inventory"
            tar -C "${SRC_DIR}" -cvzf ${target} escape-inventory
            rm "${SRC_DIR}/escape-inventory"
        else
            echo "File $target already exists"
        fi
        if [ "$INPUT_do_upload" = "1" ] ; then 
            gcs_target="gs://$INPUT_bucket/escape-inventory/$INPUT_inventory_version/$target"
            echo "Copying to $gcs_target"
            gsutil cp "$target" "$gcs_target"
            echo "Setting ACL on $gcs_target"
            gsutil acl ch -u AllUsers:R "$gcs_target"
        fi
    done
done

rm service_account.json
