/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./frontend/**/*.{ts,tsx}",
	],
	plugins: [
		require('@tailwindcss/forms'),
	],
}
