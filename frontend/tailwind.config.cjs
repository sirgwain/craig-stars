module.exports = {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		fontSize: {
			sm: ['12px', '18px'],
			base: ['14px', '20px'],
			lg: ['16px', '22px'],
			xl: ['18px', '24px']
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
