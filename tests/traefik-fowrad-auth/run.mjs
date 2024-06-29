#!/usr/bin/env -S npx zx
import 'zx/globals'
import * as zx from 'zx'

$.verbose = true

const ERR_ALREADY_EXIST = 6
const ERR_CURL_NOT_200 = 22

/**
 * 
 * @param {zx.ProcessPromise} s
 * @param {number[]} codes 
 */
async function ignoreIf(s, codes) {
	const p = await s.nothrow()
	if ([0, ...codes].includes(p.exitCode)) {
		return p
	}

	process.exit(p.exitCode)
}

/**
 * 
 * @param {string} message 
 */
function title(message) {
	echo('')
	echo(chalk.cyanBright.underline(message))
}


/**
 * 
 * @param {string} message 
 */
function sub(message) {
	echo(chalk.blackBright(message))
}

const header = 'Host: test-traefik-forward-auth'
const header_mw = header + '-mw'
const target = 'http://k3s/ping'

title('Apply manifest')
await $`kubectl apply -f ${path.join(import.meta.dirname, 'manifest.yaml')}`
await sleep(1000)

title('Test if traefik reachable')
await $`curl --fail -H ${header} ${target}`

title('Test if the middleware with forwardAuth applied (it should fail)')
if ((await ignoreIf($`curl --fail -H ${header_mw} ${target}`, [ERR_CURL_NOT_200])).exitCode === 0) {
	echo(chalk.red('it should fail'))
	process.exit(1)
}

title('Get token')
sub('Note that AlreadyExists is an expected result.')
await ignoreIf($`hr create user horus`, [ERR_ALREADY_EXIST])
const token = (await $`hr --as=horus create access-token`).stdout.trim()
await $`curl --fail -H ${header_mw} -H ${'Authorization: Bearer ' + token} ${target}`
