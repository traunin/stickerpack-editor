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
      allowedHosts: [
        env.VITE_DOMAIN,
      ],
      watch: {
        usePolling: true,
        interval: 100,
      },
      hmr: {
        host: env.VITE_DOMAIN,
        protocol: 'wss',
        clientPort: 443,
      },
    },

    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
  }
})
