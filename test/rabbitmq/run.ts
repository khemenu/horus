import 'zx/globals'
import mqtt from 'mqtt'

$.quiet = true

test('connection will require a credential', () => new Promise<void>((done, fail) => {
	const client = mqtt.connect('mqtt://rabbitmq:1883/')
	client.on('error', err => {
		if (!(err instanceof mqtt.ErrorWithReasonCode)) {
			fail('unexpected type of error')
			return
		}

		// Code 4: Connection refused since client-identifier is not allowed by the server.
		expect(err.code).to.eq(4)
		done()
	})
	client.on('connect', () => {
		fail(new Error('connection must be refused'))
	})
}))

test('connect with credential', () => new Promise<void>((done, fail) => {
	const client = mqtt.connect('mqtt://rabbitmq:1883/', {
		username: 'horus',
		password: 'horus_pass',
	})
	client.on('error', err => {
		fail(err)
	})
	client.on('connect', () => {
		done()
	})
}))
