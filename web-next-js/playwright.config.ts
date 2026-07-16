import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  // スナップショット（ベースライン画像）の保存先
  snapshotDir: "./e2e/__snapshots__",
  // スナップショットのファイル名テンプレート: {testFilePath}/{testName}.png
  snapshotPathTemplate: "{snapshotDir}/{testFilePath}/{testName}{ext}",

  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: [["html", { outputFolder: "playwright-report" }]],

  use: {
    // 画面幅固定
    viewport: { width: 1280, height: 720 },
    baseURL: "http://localhost:3000",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
    // スクリーンショット比較の許容差（2%）
    toHaveScreenshot: {
      maxDiffPixelRatio: 0.02,
      // アニメーション無効化でフレの少ない比較
      animations: "disabled",
    },
  },

  projects: [
    {
      name: "chromium",
      use: {
        channel: "chromium",
      },
    },
  ],

  // テスト実行前に Next.js dev サーバーを自動起動
  webServer: {
    command: "npm run dev",
    url: "http://localhost:3000",
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});
