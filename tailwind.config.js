/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/view/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
