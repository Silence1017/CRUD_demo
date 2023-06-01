import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    /* 是否开启 $ref */
    vue({
      refTransform: true,
      reactivityTransform: true
    }),
  ],
  server: {
    proxy: {
      '/user': {
        target: "http://127.0.0.1:3001/"
      }
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  base: './',  // 打包相对路径
})
