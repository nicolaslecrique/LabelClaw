import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  use: {
    baseURL: 'http://127.0.0.1:4173',
    trace: 'on-first-retry',
  },
  webServer: [
    {
      command: 'node tests/e2e-server.mjs',
      url: 'http://127.0.0.1:8080/api/health',
      reuseExistingServer: !process.env.CI,
      cwd: '.',
    },
    {
      command: 'pnpm vite --host 127.0.0.1 --port 4173 --mode e2e',
      url: 'http://127.0.0.1:4173',
      reuseExistingServer: !process.env.CI,
      cwd: '.',
    },
  ],
});

