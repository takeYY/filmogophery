use fred::prelude::*;
use serde::{de::DeserializeOwned, Serialize};
use tracing::warn;

use crate::config::RedisConfig;

// ─── Redis クライアント ───────────────────────────────────────

pub struct RedisClient {
    pool: Pool,
}

impl RedisClient {
    /// Redis に接続してクライアントを返す。
    /// 接続失敗は呼び出し元で `None` に変換して AppState に持たせる。
    pub async fn connect(config: &RedisConfig) -> Result<Self, anyhow::Error> {
        let url = if config.password.is_empty() {
            format!("redis://{}:{}/{}", config.host, config.port, config.db)
        } else {
            format!(
                "redis://:{}@{}:{}/{}",
                config.password, config.host, config.port, config.db
            )
        };

        let fred_config = Config::from_url(&url)?;
        let pool = Builder::from_config(fred_config).build_pool(4)?;
        pool.init().await?;

        Ok(Self { pool })
    }

    /// キャッシュを取得する。
    /// キーが存在しない・エラー・デシリアライズ失敗はすべて `None` を返す。
    pub async fn get<T: DeserializeOwned>(&self, key: &str) -> Option<T> {
        let raw: Option<String> = match self.pool.get(key).await {
            Ok(v) => v,
            Err(e) => {
                warn!("redis GET failed (key={key}): {e}");
                return None;
            }
        };

        let json = raw?;
        match serde_json::from_str::<T>(&json) {
            Ok(v) => Some(v),
            Err(e) => {
                warn!("redis deserialize failed (key={key}): {e}");
                None
            }
        }
    }

    /// キャッシュを保存する（TTL: 秒）。
    /// 保存失敗はログのみ — サービスを停止させない。
    pub async fn set<T: Serialize>(&self, key: &str, value: &T, ttl_secs: u64) {
        let json = match serde_json::to_string(value) {
            Ok(s) => s,
            Err(e) => {
                warn!("redis serialize failed (key={key}): {e}");
                return;
            }
        };

        let expire = Expiration::EX(ttl_secs as i64);
        if let Err(e) = self
            .pool
            .set::<(), _, _>(key, json, Some(expire), None, false)
            .await
        {
            warn!("redis SET failed (key={key}): {e}");
        }
    }
}
