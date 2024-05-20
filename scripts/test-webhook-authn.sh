#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset
# set -o xtrace

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)" # Directory where this script exists.
__root="$(cd "$(dirname "${__dir}")" && pwd)"         # Root directory of project.



kubectl config set-context --current --namespace=horus
if [ $(kubectl get pods --no-headers | tee /dev/tty | wc -l) -ne 1 ]; then
	echo expected that there are only 1 horus pod.
	exit 1
fi

HORUS=$(kubectl get pods -o jsonpath='{.items[0].metadata.name}')
kubectl exec "pods/${HORUS}" -- hr --conf ./conf/horus.yaml create user horus &> /dev/null || true
TOKEN=$(kubectl exec "pods/${HORUS}" -- hr --conf ./conf/horus.yaml --no-log create access-token for horus)
echo "\"${TOKEN}\""

kubectl \
	--kubeconfig="${__dir}/test-webhook-authn-config.yaml" \
	--token="${TOKEN}" \
	cluster-info
