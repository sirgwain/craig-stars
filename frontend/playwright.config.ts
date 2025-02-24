/// <reference types="node"/>
import { defineConfig } from '@playwright/test';

const command = 'npm run build && npm run preview';
const port = 4173;
// const command = 'npm run dev';
// const port = 5173;

export default defineConfig({
	retries: process.env.CI ? 2 : 0, // set to 2 when running on CI

	reporter: process.env.CI
		? [
				['github'],
				['junit', { outputFile: 'test-results.json' }],
				['html', { outputFolder: 'playwright-report', open: 'never' }]
			]
		: [['list'], ['html', { outputFolder: 'playwright-report', open: 'on-failure' }]],
	webServer: [
		{
			cwd: '../',
			command: 'go run main.go serve --test-mode',
			port: 8080,
			timeout: 60000,
			reuseExistingServer: !process.env.CI
		},
		{
			// switch commands for debugging with hot reloading
			// command: 'npm run dev',
			command: command,
			port: port
		}
	],
	expect: {
		timeout: 10000
	},
	use: {
		baseURL: `http://localhost:${port}`,
		trace: 'on-first-retry' // record traces on first retry of each test
	},

	testDir: 'e2e'
});
