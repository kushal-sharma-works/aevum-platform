/** @type {import('tailwindcss').Config} */
export default {
	darkMode: 'class',
	content: ['./index.html', './src/**/*.{vue,ts}'],
	corePlugins: {
		preflight: false
	},
	theme: {
		extend: {}
	},
	plugins: []
}
