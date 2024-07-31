import 'zx/globals'
import type { GlobalSetupContext } from 'vitest/node'

const ERR_ALREADY_EXIST = 6

export default async function setup({ }: GlobalSetupContext) {
	$.quiet = true

	// Note that it changes cwd of current process.
	cd(path.join(import.meta.dirname, '../'))

	const data = await fs.readFile('./horus.yaml', 'utf8')
	const conf = YAML.parse(data)
	const host = `localhost:${conf.http.port}`

	try {
		const ping = await fetch(`http://${host}/ping`)
		if (ping.ok) {
			console.log(`testing on a server already running in ${host}`)
			return
		}
	} catch (e) {
		if (
			![
				'UND_ERR_SOCKET', // Ill-closed port?
				'ECONNREFUSED', // Closed port.
			].includes(e.cause.code)) {
			throw e
		}
	}

	// `hr` must be installed.
	await $`go install ./cmd/hr`
	await $`hr init`
	const p = await $`hr create user horus`.nothrow()
	if (![0, ERR_ALREADY_EXIST].includes(p.exitCode)) {
		throw new Error('hr create user horus', { cause: p })
	}

	await $`echo ${'horus_pass'}`.pipe($`hr --as horus set password`)

	const horus = within(() => {
		return $`go run ./cmd/horus`.quiet(false).nothrow()
	})

	await retry(10, '0.5s', () => $`curl http://${host}/ping`)

	horus.quiet()
	return async () => horus.kill()
}
