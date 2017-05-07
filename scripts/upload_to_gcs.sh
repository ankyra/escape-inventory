#!/bin/bash

set -euf -o pipefail

tarball="escape-registry-v$INPUT_version.tgz"
target="gs://$INPUT_bucket/$tarball"

echo "Packing $tarball"
tar -cvzf "$tarball" escape-registry

echo "Uploading $INPUT_version to ${target}"
gsutil cp "$tarball" "$target"

echo "Making archive world readable"
gsutil acl ch -u AllUsers:R "$target"

