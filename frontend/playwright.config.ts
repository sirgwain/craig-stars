import type { PlaywrightTestConfig } from '@playwright/test';

const config: PlaywrightTestConfig = {
	webServer: {
		command: 'yarn run build && yarn run preview',
		port: 3000
	},
	testDir: 'tests'
};

export default config;
