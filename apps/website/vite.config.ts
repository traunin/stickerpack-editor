import path from 'path'
import process from 'process'
import vue from '@vitejs/plugin-vue'
import { defineConfig, loadEnv } from 'vite'
import svgLoader from 'vite-svg-loader'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())

  return {
    plugins: [
      vue(),
      svgLoader(),
    ],

    server: {
      port: 80,
      host: env.VITE_DOMAIN,
    },

    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
  }
})
