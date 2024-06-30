import { defineConfig } from 'vitest/config'

export default defineConfig({
	test: {
		root: './test/',
		include: ['*/**/*.ts'],
		globalSetup: [
			'setup.ts'
		],
		globals: true
	},
})
