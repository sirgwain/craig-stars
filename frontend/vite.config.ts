import { sveltekit } from '@sveltejs/kit/vite';
import type { UserConfig } from 'vite';

const config: UserConfig = {
	plugins: [sveltekit()],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
		exclude: ['node_modules', 'node_modules.nosync'],
		environment: 'jsdom',
		setupFiles: ['tests/setup.ts']
	},
	server: {
		fs: {
			allow: ['node_modules.nosync']
		},
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	},
	optimizeDeps: {
		include: ['fuzzy']
	}
};

export default config;
