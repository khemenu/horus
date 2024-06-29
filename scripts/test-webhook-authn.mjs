#!/usr/bin/env -S npx zx
import 'zx/globals'

const ERR_ALREADY_EXIST = 6

await $`kubectl config set-context --current --namespace=horus`
let p = await $`kubectl get pods --no-headers`
if (p.lines().length !== 1) {
	throw new Error(`expected that there are only 1 pod for Horus but was ${p.lines().length}`)
}

const horus_name = (await $`kubectl get pods -o jsonpath='{.items[0].metadata.name}'`).stdout
const pods_horus = `pods/${horus_name}`
echo(`test on ${pods_horus}`)

p = await $`
	kubectl exec ${pods_horus} -- \
		hr --conf ./conf/horus.yaml \
			create user horus`
	.nothrow()
if (![0, ERR_ALREADY_EXIST].includes(p.exitCode)) {
	process.exit(p.exitCode)
}

const token = (await $`
	kubectl exec ${pods_horus} -- \
		hr --conf ./conf/horus.yaml \
			--as horus \
			create access-token`)
	.stdout
	.trim()
echo(`token: ${token}`)

echo(`==========`)
echo(`test start: kubectl version`)
await $`
	kubectl \
		--kubeconfig=${path.join(__dirname, "test-webhook-authn.yaml")} \
		--token=${token} \
		version`
echo(`test passes`)
