// TODO: Redis クライアントの初期化
// fred クレートを使用予定（v10以降は tokio-native-tls または tokio-rustls featureを指定）
// 参考: https://docs.rs/fred/latest/fred/

use crate::config::RedisConfig;

pub struct RedisClient;

impl RedisClient {
    pub async fn connect(_config: &RedisConfig) -> Result<Self, anyhow::Error> {
        // TODO: fred::Client を使って接続を確立する
        Ok(RedisClient)
    }
}
