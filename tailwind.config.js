module.exports = {
	content: [
		"internal/api/template/**/*.html",
		"internal/web/template/**/*.html",
	],
	theme: {
		extend: {},
	},
	plugins: [
		require('@tailwindcss/aspect-ratio'),
		require('@tailwindcss/forms'),
	],
}
