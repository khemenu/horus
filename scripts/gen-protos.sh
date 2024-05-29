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

patch --forward \
	--reject-file - \
	"${PROTO_ROOT}/khepri/horus/horus.proto" \
	< "${__root}/scripts/horus.proto.patch"

protoc \
	--proto_path="${PROTO_ROOT}" \
	\
	--go_out="${__root}" \
	--go_opt=module="${MODULE_NAME}" \
	\
	--go-grpc_out="${__root}" \
	--go-grpc_opt=module="${MODULE_NAME}" \
	\
	--entgrpc_out="${__root}/server/bare" \
	--entgrpc_opt=module="${MODULE_NAME}" \
	--entgrpc_opt=package="${MODULE_NAME}/server/bare" \
	--entgrpc_opt=schema_path="${__root}/schema" \
	--entgrpc_opt=entity_package="${MODULE_NAME}/ent" \
	\
	"${PROTO_ROOT}"/**/*.proto
