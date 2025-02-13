/// <reference types="node"/>
import { defineConfig } from '@playwright/test';

export default defineConfig({
	reporter: [
		process.env.CI ? ['github', ['junit', { outputFile: 'test-results.json' }]] : ['list']
	],
	webServer: {
		command: 'npm run build && npm run preview',
		port: 4173
	},

	testDir: 'e2e'
});
