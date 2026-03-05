import { defineConfig } from 'vite'
import { svelte, vitePreprocess } from '@sveltejs/vite-plugin-svelte'
import { execSync } from 'child_process'

const commitHash = execSync('git rev-parse --short HEAD').toString().trim()

export default defineConfig({
  plugins: [svelte({ preprocess: vitePreprocess() })],
  define: {
    __APP_VERSION__: JSON.stringify(commitHash),
  },
  server: {
    proxy: {
      '/ws': {
        target: 'http://localhost:8080',
        ws: true,
      },
    },
  },
})
