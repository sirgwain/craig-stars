/// <reference types="vitest" />
import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';
import { fileURLToPath } from 'node:url';
import { readFileSync } from 'node:fs';
import { svelteTesting } from '@testing-library/svelte/vite';

const file = fileURLToPath(new URL('package.json', import.meta.url));
const json = readFileSync(file, 'utf8');
const pkg = JSON.parse(json);

export default defineConfig(({ mode }) => ({
	test: {
		reporters: process.env.CI ? ['junit', 'github-actions'] : 'default',
		include: ['src/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		setupFiles: ['e2e/setup.ts']
	},
	resolve: {
		conditions: mode === 'test' ? ['browser'] : []
	},
	plugins: [sveltekit(), svelteTesting()],

	define: {
		PKG: pkg
	},

	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	},
	preview: {
		// keep this the same as npm run dev to make switching modes when testing easier
		port: 5173
	},
	optimizeDeps: {
		include: ['fuzzy']
	},
	assetsInclude: ['**/*.wasm']
}));
