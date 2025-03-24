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
		reporters: ['junit', process.env.CI ? 'github-actions' : 'default'],
		include: ['src/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		outputFile: '../tmp/test-results/vitest-report.xml'
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
		},
		strictPort: true,
		open: true
	},
	optimizeDeps: {
		include: ['fuzzy']
	},
	assetsInclude: ['**/*.wasm']
}));
