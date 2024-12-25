import typography from '@tailwindcss/typography';
import daisy from 'daisyui';
import { business, emerald } from 'daisyui/src/theming/themes';
import type { Config } from 'tailwindcss';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		fontSize: {
			sm: '10px',
			base: '12px',
			lg: '14px',
			xl: '18px',
			'2xl': '20px',
			'3xl': '30px',
			'4xl': '40px',
			'5xl': '48px'
		},

		extend: {
			gridTemplateColumns: {
				// 2 column grid with an auto size label and a max value
				'label-value': 'auto minmax(0, 1fr)'
			},
			colors: {
				gauge: 'var(--gauge)'
			}
		}
	},

	plugins: [typography, daisy],

	darkMode: 'selector',

	daisyui: {
		themes: [
			{
				business: {
					...business,
					'base-100': '#252525',
					'base-200': '#212121',
					'base-300': '#151515',
					'--gauge': '#151515'
				},
				emerald: {
					...emerald,
					primary: '#4D9A69',
					'base-200': '#C3C3C3', // win31!
					'--gauge': '#555555'
				}
			}
		],
		darkTheme: 'business'
	}
} as Config;
