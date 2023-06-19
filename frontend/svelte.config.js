import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: [
		preprocess({
			postcss: true
		})
	],
	kit: {
		adapter: adapter({
			fallback: 'index.html'
		}),
		prerender: {
			enabled: false,
			onError: 'continue'
		},
		vite: {
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
			}
		}
	}
};

export default config;
