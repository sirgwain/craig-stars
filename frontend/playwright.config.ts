/// <reference types="node"/>
import { defineConfig } from '@playwright/test';

const command = 'npm run build && npm run preview';
const port = 4173;
// const command = 'npm run dev';
// const port = 5173;

export default defineConfig({
	retries: process.env.CI ? 2 : 0, // set to 2 when running on CI
	// single threaded or we need to use a non memory db
	workers: 1,

	reporter: [
		process.env.CI ? ['github'] : ['list'],
		[
			'html',
			{
				outputFolder: '../tmp/test-results/playwright-reports',
				open: process.env.CI ? 'never' : 'on-failure'
			}
		],

		// Change these if we ever change github action's tmpdir folder
		['junit', { outputFile: '../tmp/test-results/playwright-reports/playwright-report.xml' }]
	],
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
			port: port,
			reuseExistingServer: !process.env.CI
		}
	],
	expect: {
		timeout: 15000
	},
	use: {
		baseURL: `http://localhost:${port}`,
		headless: true, // set to false to see cool popup windows
		trace: 'on-first-retry' // record traces on first retry of each test
		// video: 'on' // turn on for cool video recordings
	},

	testDir: 'e2e'
});
