/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/html/**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
