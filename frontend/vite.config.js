import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    host: '127.0.0.1',
    proxy: {
      '/.ory': {
        target: 'http://127.0.0.1:4433',
        changeOrigin: true
      },
      '/self-service': {
        target: 'http://127.0.0.1:4433',
        changeOrigin: true
      },
      '/sessions': {
        target: 'http://127.0.0.1:4433',
        changeOrigin: true
      },
      '/identity': {
        target: 'http://127.0.0.1:4433',
        changeOrigin: true
      },
      '/api': {
        target: 'http://127.0.0.1:8888',
        changeOrigin: true
      }
    }
  }
})