import { defineConfig } from 'vite';

export default defineConfig({
  server: {
    host: '127.0.0.1',
    port: Number(process.env.UI_DEV_PORT || 5173),
    strictPort: true,
    proxy: {
      '/api/runs': {
        target: process.env.API_UPSTREAM || 'http://127.0.0.1:8080',
        rewrite: () => '/app/cicd/api/runs',
      },
    },
  },
});
