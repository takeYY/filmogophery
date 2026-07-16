use tracing_subscriber::{EnvFilter, fmt};

use crate::config::LogConfig;

/// ロガーを初期化する
/// LOG_LEVEL 環境変数が優先される（例: "debug", "info"）
pub fn init(config: &LogConfig) {
    let filter = EnvFilter::try_from_default_env()
        .unwrap_or_else(|_| EnvFilter::new(&config.level));

    fmt::Subscriber::builder()
        .with_env_filter(filter)
        .with_target(true)
        .init();
}
