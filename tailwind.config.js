module.exports = {
	content: [
		"internal/api/templates/*.html",
		"internal/web/templates/*.html",
	],
	theme: {
		extend: {},
	},
	plugins: [
		require('@tailwindcss/aspect-ratio'),
		require('@tailwindcss/forms'),
	],
}
