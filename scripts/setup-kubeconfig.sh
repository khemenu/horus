#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.



DST="${HOME}/.kube"
mkdir -p "${DST}"
sed 's/127.0.0.1/k3s/g' /etc/kubeconfig/kubeconfig.yaml > "${DST}/config"

kubectl cluster-info
kubectl create namespace horus 2> /dev/null || true
kubectl config set-context --current --namespace=horus
