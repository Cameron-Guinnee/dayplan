import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/tasks': 'http://localhost:8080',
      '/time-blocks': 'http://localhost:8080',
      '/schedule': 'http://localhost:8080',
    },
  },
})
