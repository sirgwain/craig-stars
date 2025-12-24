import prettier from 'eslint-config-prettier';
import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';
import ts from 'typescript-eslint';

export default ts.config(
	js.configs.recommended,
	...ts.configs.recommended,
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
		files: ['**/*.svelte'],

		languageOptions: {
			parserOptions: {
				parser: ts.parser,
				project: ['./tsconfig.json'],
				tsconfigRootDir: process.cwd(),
				extraFileExtensions: ['.svelte']
			}
		},

		rules: {
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', caughtErrors: 'all', caughtErrorsIgnorePattern: '^_' }
			],
			'@typescript-eslint/no-unnecessary-condition': 'error',
			'no-undef': 'off'
		}
	},
	{
		ignores: [
			'src/lib/wasm/wasm_exec.js',
			'src/lib/protogen',
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
