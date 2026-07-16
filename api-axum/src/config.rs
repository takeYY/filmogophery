use serde::Deserialize;

/// envy はネストした struct に対応していないため、
/// 環境変数をフラットな構造体で受け取ってから Config に変換する
#[derive(Debug, Deserialize)]
struct Env {
    // Server
    server_port: u16,

    // Logger
    log_level: String,

    // Writer DB
    writer_db_host: String,
    writer_db_name: String,
    writer_db_user: String,
    writer_db_pwd: String,
    writer_db_core_count: u32,

    // Reader DB
    reader_db_host: String,
    reader_db_name: String,
    reader_db_user: String,
    reader_db_pwd: String,
    reader_db_core_count: u32,

    // Redis
    redis_host: String,
    redis_port: u16,
    redis_password: String,
    redis_db: u8,

    // TMDB
    tmdb_access_token: String,

    // JWT
    jwt_secret: String,
}

/// アプリケーション全体の設定
#[derive(Debug)]
pub struct Config {
    pub server: ServerConfig,
    pub database: DatabaseConfig,
    pub redis: RedisConfig,
    pub log: LogConfig,
    pub jwt: JwtConfig,
    pub tmdb: TmdbConfig,
}

#[derive(Debug)]
pub struct ServerConfig {
    pub port: u16,
}

#[derive(Debug)]
pub struct DatabaseConfig {
    pub writer_host: String,
    pub writer_name: String,
    pub writer_user: String,
    pub writer_password: String,
    pub writer_core_count: u32,

    pub reader_host: String,
    pub reader_name: String,
    pub reader_user: String,
    pub reader_password: String,
    pub reader_core_count: u32,
}

#[derive(Debug)]
pub struct RedisConfig {
    pub host: String,
    pub port: u16,
    pub password: String,
    pub db: u8,
}

#[derive(Debug)]
pub struct LogConfig {
    pub level: String,
}

#[derive(Debug)]
pub struct JwtConfig {
    pub secret: String,
}

#[derive(Debug)]
pub struct TmdbConfig {
    pub access_token: String,
}

impl Config {
    pub fn from_env() -> Result<Self, envy::Error> {
        let env = envy::from_env::<Env>()?;

        Ok(Config {
            server: ServerConfig {
                port: env.server_port,
            },
            database: DatabaseConfig {
                writer_host: env.writer_db_host,
                writer_name: env.writer_db_name,
                writer_user: env.writer_db_user,
                writer_password: env.writer_db_pwd,
                writer_core_count: env.writer_db_core_count,
                reader_host: env.reader_db_host,
                reader_name: env.reader_db_name,
                reader_user: env.reader_db_user,
                reader_password: env.reader_db_pwd,
                reader_core_count: env.reader_db_core_count,
            },
            redis: RedisConfig {
                host: env.redis_host,
                port: env.redis_port,
                password: env.redis_password,
                db: env.redis_db,
            },
            log: LogConfig {
                level: env.log_level,
            },
            jwt: JwtConfig {
                secret: env.jwt_secret,
            },
            tmdb: TmdbConfig {
                access_token: env.tmdb_access_token,
            },
        })
    }
}
