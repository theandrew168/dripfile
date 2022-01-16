module.exports = {
	content: [
		"internal/api/templates/*.tmpl",
		"internal/web/templates/*.tmpl",
	],
	theme: {
		extend: {},
	},
	plugins: [
		require('@tailwindcss/aspect-ratio'),
		require('@tailwindcss/forms'),
	],
}
