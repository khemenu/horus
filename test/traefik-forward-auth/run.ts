import 'zx/globals'

const ERR_CURL_NOT_200 = 22

const name = 'test-traefik-forward-auth'
const header = `Host: ${name}`
const headerMw = header + '-mw'
const target = 'http://k3s/ping'

$.quiet = true

beforeAll(async () => {
	const manifest = path.join(import.meta.dirname, 'manifest.yaml')
	await $`kubectl apply -f ${manifest}`
	await retry(10, '0.3s', () => $`curl --fail -H ${header} http://k3s/api/http/middlewares/kube-system-${name}@kubernetescrd`)

	return () => $`kubectl delete -f ${manifest}`
})

test('middleware forwards the traffic to auth', async () => {
	const p = await $`curl --fail -H ${headerMw} ${target}`.nothrow()
	expect(p.exitCode).toBe(ERR_CURL_NOT_200)
	expect(p.stderr).contains('401')
})

test('request with invalid Authorization format', async () => {
	const p = await $`curl --fail -H ${headerMw} -H "Authorization: invalid" ${target}`.nothrow()
	expect(p.exitCode).toBe(ERR_CURL_NOT_200)
	expect(p.stderr).contains('401')
})

test('request with invalid token', async () => {
	// Well-formed but invalid token.
	const token = "a".repeat(16) + 'bbbb'

	const p = await $`curl --fail -H ${headerMw} -H "Authorization: Bearer ${token}" ${target}`.nothrow()
	console.log(p.stderr)
	expect(p.exitCode).toBe(ERR_CURL_NOT_200)
	expect(p.stderr).contains('401')
})

test('request with valid token', async () => {
	const token = (await $`hr --as=horus create access-token`).stdout.trim()
	await $`curl --fail -H ${headerMw} -H ${'Authorization: Bearer ' + token} ${target}`
})
