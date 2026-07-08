/// <reference types="vitest/config" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Dev: proxy para o gateway Go local. Build: sai em backend-gateway/web/dist
// para ser embutido no binário via go:embed.
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:8090',
      '/auth': 'http://localhost:8090',
    },
  },
  build: {
    outDir: '../backend-gateway/web/dist',
    emptyOutDir: true,
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./vitest.setup.ts'],
  },
})
