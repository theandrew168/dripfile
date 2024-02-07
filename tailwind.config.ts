import type { Config } from 'tailwindcss';
import forms from '@tailwindcss/forms';

export default {
	content: [
		"./index.html",
		"./frontend/**/*.{ts,tsx}",
	],
	plugins: [
		forms,
	],
} satisfies Config;
