import 'zx/globals'

$.quiet = true

test.fails('request with invalid token', async () => {
	await $`kubectl --kubeconfig=${path.join(import.meta.dirname, './config.yaml')} --token=invalid version`
})

test('request with valid token', async () => {
	const token = (await $`hr --as=horus create access-token`).stdout.trim()
	await $`kubectl --kubeconfig=${path.join(import.meta.dirname, './config.yaml')} --token=${token} version`
})
