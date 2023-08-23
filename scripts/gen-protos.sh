#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.


cd "${__root}/pb"
rm *.pb.go
protoc \
	--proto_path=${__root}/pb \
	--plugin=protoc-gen-go=$(which protoc-gen-go) \
	--plugin=protoc-gen-go-grpc=$(which protoc-gen-go-grpc) \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	*.proto
