module.exports = {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		fontSize: {
			sm: ['10px'],
			base: ['12px'],
			lg: ['14px'],
			xl: ['18px'],
			'2xl': ['20px'],
			'3xl': ['30px'],
			'4xl': ['40px'],
			'5xl': ['48px']
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

	plugins: [require('@tailwindcss/typography'), require('daisyui')],

	darkMode: ['class'],

	daisyui: {
		themes: [
			{
				business: {
					...require('daisyui/src/theming/themes')['[data-theme=business]'],
					'base-100': '#252525',
					'base-200': '#212121',
					'base-300': '#151515',
					'--gauge': '#151515'
				},
				emerald: {
					...require('daisyui/src/theming/themes')['[data-theme=emerald]'],
					primary: '#4D9A69',
					'base-200': '#C3C3C3', // win31!
					'--gauge': '#555555'
				}
			}
		],
		darkTheme: 'business'
	}
};
