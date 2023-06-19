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

	plugins: [require('@tailwindcss/forms'), require('@tailwindcss/typography'), require('daisyui')],

	darkMode: 'class',

	daisyui: {
		themes: ['emerald', 'business'],
		darkTheme: 'business'
	}
};
