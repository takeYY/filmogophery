# api-axum

Rust / Axum による REST API 実装。Echo（Go）・Hono（TypeScript）と同一の機能を Rust で実装することで、言語・フレームワーク間の実装差異を学ぶことを目的としています。

## Tech Stack

| カテゴリ           | ライブラリ                                                                               |
| ------------------ | ---------------------------------------------------------------------------------------- |
| Web フレームワーク | [axum](https://github.com/tokio-rs/axum) 0.8                                             |
| 非同期ランタイム   | [tokio](https://tokio.rs) 1                                                              |
| ミドルウェア       | [tower-http](https://github.com/tower-rs/tower-http)（CORS、トレース）                   |
| DB                 | MySQL / [sqlx](https://github.com/launchbadge/sqlx) 0.8                                  |
| 認証               | JWT HS256 / [jsonwebtoken](https://github.com/Keats/jsonwebtoken)                        |
| パスワード         | [bcrypt](https://github.com/Keats/rust-bcrypt)                                           |
| バリデーション     | [validator](https://github.com/Keats/validator)                                          |
| エラーハンドリング | [thiserror](https://github.com/dtolnay/thiserror)                                        |
| 設定               | [envy](https://github.com/softprops/envy) + [dotenvy](https://github.com/allan2/dotenvy) |
| ロガー             | [tracing](https://github.com/tokio-rs/tracing) + tracing-subscriber                      |
| テスト             | [axum-test](https://github.com/JoshuaColell/axum-test)                                   |

## Directory Structure

```
api-axum/
├── Cargo.toml               # 依存クレート定義
├── rust-toolchain.toml      # Rust 1.88.0 を固定
├── .env                     # 環境変数（git管理外）
├── .example.env             # 環境変数のサンプル
└── src/
    ├── main.rs              # エントリーポイント（設定読み込み・DB接続・サーバー起動）
    ├── config.rs            # 環境変数 → Config 構造体への変換
    ├── app/
    │   ├── mod.rs
    │   ├── router.rs        # AppState定義・ルーター組み立て
    │   ├── responses.rs     # AppError（thiserror）+ IntoResponse実装
    │   └── features/        # 機能ごとのモジュール
    │       ├── health/
    │       │   ├── mod.rs
    │       │   ├── handler.rs      # HTTPハンドラ・ルート定義・テスト
    │       │   └── use_case.rs     # ビジネスロジック
    │       ├── user/
    │       │   ├── mod.rs
    │       │   ├── handler.rs      # HTTPハンドラ・ルート定義
    │       │   ├── use_case.rs     # ビジネスロジック（repository trait を呼び出す）
    │       │   └── repository.rs   # DB アクセス（trait 定義 + MySQL 実装）
    │       ├── auth/
    │       ├── movie/
    │       ├── review/
    │       ├── watchlist/
    │       ├── trending/
    │       ├── search/
    │       └── master/
    └── pkg/                 # 共通ユーティリティ
        ├── db.rs            # MySQL 接続プール（sqlx）
        ├── jwt.rs           # JWT 生成・検証
        ├── logger.rs        # tracing 初期化
        ├── redis.rs         # Redis クライアント
        └── middleware/
            └── auth.rs      # Bearer JWT 認証ミドルウェア + AuthUser エクストラクタ
```

### アーキテクチャ

各 feature は 3 層に分離しています。Echo 実装の Handler / Interactor / Repository パターンに対応します。

```
Request
  └─ handler.rs      # Bind / Validate → use_case 呼び出し → Response
       └─ use_case.rs    # ビジネスロジック（repository trait に依存）
            └─ repository.rs  # trait 定義 + MySQL 実装（sqlx による SQL 発行）
```

`repository.rs` は trait（`UserRepository` など）と MySQL 実装（`MySqlUserRepository` など）を同一ファイルに置きます。将来テストでモックに差し替える場合は trait を実装した別の struct を用意します。

> `health` など DB アクセスが不要な feature は `repository.rs` を持ちません。

## How to use it

### Docker で起動（推奨）

```sh
# ルートの Makefile から起動
make up TARGET_COMPOSE=compose.axum.yml

# バックグラウンドで起動
make up_d TARGET_COMPOSE=compose.axum.yml

# 停止
make stop TARGET_COMPOSE=compose.axum.yml
```

### ローカルで起動

```sh
# 依存クレートのビルド
cargo build

# 開発サーバー起動（ホットリロードなし）
cargo run

# ホットリロードあり（cargo-watch が必要）
cargo install cargo-watch --locked
cargo watch -x run
```

### テスト実行

```sh
# ユニットテスト
cargo test

# Docker 経由でテスト
make test_axum
```

## API Endpoints

実装状況は [docs/readme.md](../docs/readme.md) の Feature テーブルを参照してください。

## Notes

- Rust ツールチェーンは `rust-toolchain.toml` で 1.88.0 に固定しています。IDE のエラーが出る場合は以下を実行してください。
  ```sh
  rustup component add rust-analyzer rust-src --toolchain 1.88.0
  ```
- `AppState.db` は現在 `Option<MySqlPool>` です。DB 接続に失敗してもサーバーは起動します。各機能を実装する際は `state.db.as_ref()` で接続の有無を確認してください。
