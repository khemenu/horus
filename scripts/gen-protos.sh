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

protoc \
	--proto_path="${PROTO_ROOT}" \
	\
	--go_out="${__root}" \
	--go_opt=module="${MODULE_NAME}" \
	\
	--go-grpc_out="${__root}" \
	--go-grpc_opt=module="${MODULE_NAME}" \
	\
	--entpb_out="${__root}" \
	--entpb_opt=module="${MODULE_NAME}" \
	--entpb_opt=schema_path="${__root}/schema" \
	--entpb_opt=ent_package="${MODULE_NAME}/ent" \
	--entpb_opt=package="${MODULE_NAME}/server/bare" \
	\
	"${PROTO_ROOT}"/**/*.proto
