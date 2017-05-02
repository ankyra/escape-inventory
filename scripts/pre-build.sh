#!/bin/bash -e

set -euf -o pipefail

cat > metadata.go <<EOF
package main

const registryVersion="$INPUT_version"
EOF
