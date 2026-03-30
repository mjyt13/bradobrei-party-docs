import { defineConfig, loadEnv } from 'vite'
import react, { reactCompilerPreset } from '@vitejs/plugin-react'
import babel from '@rolldown/plugin-babel'
import { createViteProxy } from './vite.proxy'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const devHost = env.VITE_DEV_HOST || 'localhost'
  const devPort = Number(env.VITE_DEV_PORT || '5173')

  return {
    plugins: [react(), babel({ presets: [reactCompilerPreset()] })],
    server: {
      host: devHost,
      port: devPort,
      strictPort: true,
      hmr: {
        host: devHost,
      },
      proxy: createViteProxy(env),
    },
  }
})
