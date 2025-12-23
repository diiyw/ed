import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    proxy: {
      // Proxy API requests to backend
      '/api': {
        target: process.env.VITE_API_PROXY_TARGET || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
      // Proxy WebSocket connections for deployment logs
      '/ws': {
        target: process.env.VITE_WS_PROXY_TARGET || 'ws://localhost:8080',
        ws: true,
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
