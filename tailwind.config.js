/** @type {import('tailwindcss').Config} */
module.exports = {
	content: [
		"./public/index.html",
		"./frontend/**/*.{ts,tsx}",
	],
	plugins: [
		require('@tailwindcss/forms'),
	],
}
