#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

shopt -s globstar

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.



MODULE_NAME=khepri.dev/horus
PROTO_ROOT="${__root}/proto"
cd "${PROTO_ROOT}"

# There is TODO in the code.
# https://github.com/ent/contrib/blob/4ec197664a206890a44245f5c0cbcb8110d68cb5/entproto/adapter.go#L206C2-L206C62
sed -i 's/khepri.dev\/horus\/ent\/proto\/khepri\/horus/khepri.dev\/horus/g' khepri/horus/horus.proto

protoc \
	--proto_path="${PROTO_ROOT}" \
	\
	--go_out="${__root}" \
	--go_opt=module="${MODULE_NAME}" \
	\
	--go-grpc_out="${__root}" \
	--go-grpc_opt=module="${MODULE_NAME}" \
	\
	--entgrpc_out="${__root}/service/bare" \
	--entgrpc_opt=module="${MODULE_NAME}" \
	--entgrpc_opt=package="${MODULE_NAME}/service/bare" \
	--entgrpc_opt=schema_path="${__root}/schema" \
	--entgrpc_opt=entity_package="${MODULE_NAME}/ent" \
	\
	"${PROTO_ROOT}"/**/*.proto
