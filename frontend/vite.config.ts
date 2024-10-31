import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';
import { fileURLToPath } from 'node:url';
import { readFileSync } from 'node:fs';
import { svelteTesting } from '@testing-library/svelte/vite';

const file = fileURLToPath(new URL('package.json', import.meta.url));
const json = readFileSync(file, 'utf8');
const pkg = JSON.parse(json);

export default defineConfig({
	plugins: [sveltekit(), svelteTesting()],

	define: {
		PKG: pkg
	},

	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		setupFiles: ['e2e/setup.ts']
	},
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	},
	optimizeDeps: {
		include: ['fuzzy']
	},
	assetsInclude: ['**/*.wasm']
});
