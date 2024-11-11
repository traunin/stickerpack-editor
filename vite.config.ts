import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import basicSsl from '@vitejs/plugin-basic-ssl'
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    basicSsl()
  ],

  server: {
    port: 443
  },
  
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    },
  },
})
