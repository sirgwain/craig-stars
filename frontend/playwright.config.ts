/// <reference types="node"/>
import { defineConfig } from '@playwright/test';

const command = 'npm run build && npm run preview';
const port = 4173;
// const command = 'npm run dev';
// const port = 5173;

export default defineConfig({
	// Global setup to authenticate once
	globalSetup: './e2e/global.setup.ts',

	retries: process.env.CI ? 2 : 0, // set to 2 when running on CI
	// single threaded or we need to use a non memory db
	workers: 5,

	projects: [
		{
			name: 'login-tests',
			testMatch: '**/login.spec.ts',
			use: {
				baseURL: `http://localhost:${port}`,
				headless: true,
				trace: 'on-first-retry'
				// No storageState - these tests should start without authentication
			}
		},
		{
			name: 'authenticated-tests',
			testIgnore: '**/login.spec.ts',
			use: {
				baseURL: `http://localhost:${port}`,
				headless: true,
				trace: 'on-first-retry',
				storageState: '.auth/user.json' // Use cached authentication state
			}
		}
	],

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

	testDir: 'e2e'
});
