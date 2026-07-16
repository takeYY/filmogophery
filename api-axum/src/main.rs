use std::sync::Arc;

use tracing::info;

mod app;
mod config;
mod pkg;

use app::router::create_router;
use config::Config;
use pkg::db;

#[tokio::main]
async fn main() {
    // .env 読み込み（存在しない場合はスキップ）
    dotenvy::dotenv().ok();

    let config = Config::from_env().expect("Failed to load config");

    // ロガー初期化（DB接続より先に行う）
    pkg::logger::init(&config.log);

    // DB 接続
    // 接続に失敗してもサーバーは起動させる（ヘルスチェック等が機能するように）
    // TODO: 実装が進んだら None の場合はサーバーを起動しない形に変える
    let db_pool = match db::connect(&config.database).await {
        Ok(pool) => {
            info!("Successfully connected to database");
            Some(pool)
        }
        Err(e) => {
            tracing::warn!("Failed to connect to database: {}. Starting without DB.", e);
            None
        }
    };

    let config = Arc::new(config);

    let router = create_router(config.clone(), db_pool);

    let addr = format!("0.0.0.0:{}", config.server.port);
    let listener = tokio::net::TcpListener::bind(&addr)
        .await
        .expect("Failed to bind address");

    info!("Starting server on {}", addr);
    axum::serve(listener, router)
        .await
        .expect("Server error");
}
