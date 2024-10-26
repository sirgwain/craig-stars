import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations#preprocessors
	// for more information about preprocessors
	// this is because we are using some constructor accessors
	preprocess: vitePreprocess({ script: true }),

	kit: {
		// adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://svelte.dev/docs/kit/adapters for more information about adapters.
		adapter: adapter({
			fallback: 'index.html'
		}),
		prerender: { entries: [] }
	},
	// onwarn: (warning, handler) => {
	// 	if (warning.code === 'reactive_declaration_non_reactive_property') return;
	//
	// 	handler(warning);
	// }
};

export default config;
