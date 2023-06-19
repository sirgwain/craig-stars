module.exports = {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {}
	},

	plugins: [require('@tailwindcss/forms'), require('@tailwindcss/typography'), require('daisyui')],

	darkMode: 'class',

	daisyui: {
		themes: ['emerald', 'business'],
		darkTheme: 'business'
	}
};
