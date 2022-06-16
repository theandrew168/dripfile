import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
	build: {
		outDir: 'backend/web/public',
  },
	clearScreen: false,
	plugins: [svelte()],
	publicDir: 'frontend/public',
})
