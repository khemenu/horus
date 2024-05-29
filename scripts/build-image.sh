#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

shopt -s globstar

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.



TEMP=$(mktemp)
cd "${__root}"

tar -cf "${TEMP}" \
	./alias/ \
	./cmd/ \
	./ent/ \
	./internal/ \
	./log/ \
	./proto/ \
	./schema/ \
	./server/ \
	./tokens/ \
	./Dockerfile \
	./go.mod \
	./go.sum \
	./*.go

docker build -t registry:5000/khepri/horus:latest - < "${TEMP}"
docker push registry:5000/khepri/horus:latest

rm "${TEMP}"
