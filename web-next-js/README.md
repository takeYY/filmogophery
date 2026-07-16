# Filmogophery Next

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/basic-features/font-optimization) to automatically optimize and load Inter, a custom Google Font.

## E2E Testing (Playwright)

画面幅 1280x720 固定のスクリーンショット比較テストを Playwright で実施します。

### コマンド

```bash
# テスト実行（dev サーバーは自動起動）
npm run test:e2e

# ベースライン画像の初回生成・更新
npm run test:e2e:update

# UI モードで対話的にテストを確認
npm run test:e2e:ui

# HTML レポートを開く
npm run test:e2e:report
```

### ローカルで Docker を使って実行する（CI と環境を統一）

フォントレンダリングを CI と完全に一致させるには、Playwright 公式イメージ上で実行します。

```bash
# プロジェクトルートから実行
docker run --rm \
  -v $(pwd)/web-next-js:/work \
  -w /work \
  mcr.microsoft.com/playwright:v1.61.1-noble \
  bash -c "npm ci && npm run test:e2e"

# ベースライン画像の更新
docker run --rm \
  -v $(pwd)/web-next-js:/work \
  -w /work \
  mcr.microsoft.com/playwright:v1.61.1-noble \
  bash -c "npm ci && npm run test:e2e:update"
```

> **Note:** ベースライン画像は必ずこの Docker コマンドで生成・更新してからコミットしてください。ホスト OS で生成した画像は CI と一致しない場合があります。

### 特定の画面だけ実行する

```bash
# ファイル名で絞る
npm run test:e2e -- e2e/home.spec.ts

# テスト名（describe/test の文字列）で絞る
npm run test:e2e -- -g "Home page"

# 複数ファイルを指定
npm run test:e2e -- e2e/home.spec.ts e2e/login.spec.ts
```

### ベースライン画像について

- `e2e/__snapshots__/` 以下のベースライン画像は git で管理します
- `test-results/`（失敗時の差分画像）と `playwright-report/` は git 管理外です
- 初回セットアップ時や UI 変更後は `npm run test:e2e:update` でベースラインを更新してからコミットしてください
