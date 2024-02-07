/** @type {import('tailwindcss').Config} */
export default {
	content: [
		"./index.html",
		"./frontend/**/*.{ts,tsx}",
	],
	plugins: [
		require('@tailwindcss/forms'),
	],
}
