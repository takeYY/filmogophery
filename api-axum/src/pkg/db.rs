use sqlx::{MySqlPool, mysql::MySqlPoolOptions};

use crate::config::DatabaseConfig;

/// MySQL 接続プールを作成する（Writer / Reader 両対応）
/// 現在は Writer 用の単一プールを返す。
/// Reader/Writer 分離は AppState で複数プールを持つ形で拡張する。
pub async fn connect(config: &DatabaseConfig) -> Result<MySqlPool, sqlx::Error> {
    let dsn = format!(
        "mysql://{}:{}@{}/{}",
        config.writer_user,
        config.writer_password,
        config.writer_host,
        config.writer_name,
    );

    let max_connections = config.writer_core_count * 2;

    MySqlPoolOptions::new()
        .max_connections(max_connections)
        .connect(&dsn)
        .await
}
