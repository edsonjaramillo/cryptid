import tailwindcss from '@tailwindcss/vite';
import { TanStackRouterVite } from '@tanstack/router-plugin/vite';
import viteReact from '@vitejs/plugin-react';
import { defineConfig } from 'vite';
import { ENV } from './src/utils/env';

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: ENV.WEB_PORT || 3000,
  },
  plugins: [TanStackRouterVite({ autoCodeSplitting: true }), viteReact(), tailwindcss()],
});
