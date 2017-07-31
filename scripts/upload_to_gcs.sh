#!/bin/bash

set -euf -o pipefail

tarball="escape-registry-v$INPUT_version.tgz"
target="gs://$INPUT_bucket/escape-registry/$INPUT_version/$tarball"

echo "Packing $tarball"
tar -cvzf "$tarball" escape-registry

echo "$INPUT_credentials" > service_account.json
gcloud auth activate-service-account --key-file service_account.json

echo "Uploading $INPUT_version to ${target}"
gsutil cp "$tarball" "$target"

echo "Making archive world readable"
gsutil acl ch -u AllUsers:R "$target"

rm service_account.json
