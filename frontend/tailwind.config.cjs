module.exports = {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		fontSize: {
			sm: ['10px'],
			base: ['12px'],
			lg: ['14px'],
			xl: ['18px'],
			'2xl': ['20px']
		},

		extend: {}
	},

	plugins: [require('@tailwindcss/typography'), require('daisyui')],

	darkMode: 'class',

	daisyui: {
		themes: [
			{
				business: {
					...require("daisyui/src/colors/themes")["[data-theme=business]"],
					// 'base-100': '#333333'
				}
			},
			'emerald',
		],
		darkTheme: 'business'
	}
};
