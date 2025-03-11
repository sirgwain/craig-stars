import prettier from 'eslint-config-prettier';
import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';
import ts from 'typescript-eslint';
import svelteConfig from './svelte.config.js';

export default ts.config(
	js.configs.recommended,
	...ts.configs.recommended,
	//	...ts.configs.stylistic,
	...svelte.configs['flat/recommended'],
	prettier,
	...svelte.configs['flat/prettier'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			}
		}
	},
	{
		rules: {
			'no-var': 'error'
		}
	},
	{
		files: ['**/*.svelte'],

		languageOptions: {
			parserOptions: {
				parser: ts.parser,
				svelteConfig
			}
		},

		rules: {
			'svelte/require-each-key': 'off',
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', caughtErrors: 'all', caughtErrorsIgnorePattern: '^_' }
			],
			'no-undef': 'off'
		}
	},
	{
		ignores: [
			'src/lib/wasm/wasm_exec.js',
			'!.env.example',
			'.DS_Store',
			'.env.*',
			'.svelte-kit/',
			'/package',
			'build/',
			'dist/',
			'node_modules',
			'package-lock.json',
			'pnpm-lock.yaml',
			'yarn.lock'
		]
	}
);
