import prettier from 'eslint-config-prettier';
import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';
import ts from 'typescript-eslint';

export default ts.config(
	js.configs.recommended,
	...ts.configs.recommendedTypeChecked,
	...svelte.configs['flat/recommended'],
	prettier,
	...svelte.configs['flat/prettier'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			},
			parserOptions: {
				extraFileExtensions: [".svelte"],
				projectService: true,
				tsconfigRootDir: import.meta.dirname,
			},
		},
	},
	{
		rules: {
			'no-var': 'error',
			'@typescript-eslint/ban-ts-comment': [
				'error',
				{
					'ts-expect-error': {descriptionFormat: '^: .+$'},
				},
			],
			'@typescript-eslint/no-unnecessary-boolean-literal-compare': 'error',
			'@typescript-eslint/no-unnecessary-condition': 'error',
			'@typescript-eslint/consistent-type-assertions': 'warn',
			'@typescript-eslint/no-confusing-non-null-assertion': 'error'
		}
	},
	{
		files: ['**/*.svelte'],

		languageOptions: {
			parserOptions: {
				parser: ts.parser
			}
		},

		rules: {
			'@typescript-eslint/no-unused-vars': [
				'error',
				{ argsIgnorePattern: '^_', caughtErrors: 'all', caughtErrorsIgnorePattern: '^_' }
			],
			'no-undef': 'off',
			'@typescript-eslint/no-unsafe-member-access': 'off',
			'@typescript-eslint/no-unsafe-argument': 'off',
			'@typescript-eslint/no-unsafe-assignment': 'off',
			'@typescript-eslint/no-unsafe-call': 'off',
			'@typescript-eslint/no-unsafe-condition': 'off',

		}
	},
	{
		// disable type aware linting on config files (they're for config after all)
		files: ['**/*.config.{js, ts}'],
		extends: [ts.configs.disableTypeChecked],
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
