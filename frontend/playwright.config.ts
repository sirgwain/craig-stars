/// <reference types="node"/>
import { defineConfig } from '@playwright/test';

export default defineConfig({
	reporter: [
		process.env.CI ? ['github', ['junit', { outputFile: 'test-results.json' }]] : ['list']
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
			command: 'npm run build && npm run preview',
			port: 5173
		}
	],
	use: {
		baseURL: 'http://localhost:5173'
	},

	testDir: 'e2e'
});
