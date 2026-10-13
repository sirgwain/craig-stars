import prettier from 'eslint-config-prettier';
import js from '@eslint/js';
import betterTailwindcss from 'eslint-plugin-better-tailwindcss';
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
		},
		rules: {
			'svelte/no-navigation-without-resolve': 'off'
		}
	},
	{
		files: ['**/*.svelte', '**/*.svelte.ts'],

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
		// the same tailwind hints the vscode extension shows
		files: ['**/*.svelte'],
		plugins: {
			'better-tailwindcss': betterTailwindcss
		},
		settings: {
			'better-tailwindcss': {
				entryPoint: 'src/css/app.css'
			}
		},
		rules: {
			'better-tailwindcss/enforce-canonical-classes': 'error',
			'better-tailwindcss/no-conflicting-classes': 'error',
			'better-tailwindcss/no-duplicate-classes': 'error'
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
