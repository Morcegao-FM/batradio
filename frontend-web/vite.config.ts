/// <reference types="vitest/config" />
import { writeFileSync } from 'node:fs'
import { join } from 'node:path'
import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'

// O go:embed do gateway exige que web/dist exista em clones frescos; o
// .gitkeep é versionado, mas o emptyOutDir do build o apaga — recria aqui.
function keepGitkeep(outDir: string): Plugin {
  return {
    name: 'keep-gitkeep',
    closeBundle() {
      writeFileSync(join(outDir, '.gitkeep'), '')
    },
  }
}

// Dev: proxy para o gateway Go local. Build: sai em backend-gateway/web/dist
// para ser embutido no binário via go:embed.
export default defineConfig({
  plugins: [react(), keepGitkeep('../backend-gateway/web/dist')],
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
