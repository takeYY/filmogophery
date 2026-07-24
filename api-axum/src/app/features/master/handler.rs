use axum::{Json, Router, extract::State, routing::get};
use tracing::info;

use crate::app::responses::{ApiResult, AppError};
use crate::app::router::AppState;

use super::repository::MySqlMasterRepository;
use super::use_case::{self, GenreResponse, PlatformResponse};

// ─── ハンドラ ──────────────────────────────────────────────────

/// GET /v1/genres — ジャンル一覧（認証必要）
async fn get_genres(
    State(state): State<AppState>,
) -> ApiResult<Json<Vec<GenreResponse>>> {
    info!("accessed GET /v1/genres");

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let repo = MySqlMasterRepository(db);

    let genres = use_case::get_genres(&repo).await?;
    Ok(Json(genres))
}

/// GET /v1/platforms — プラットフォーム一覧（認証必要）
async fn get_platforms(
    State(state): State<AppState>,
) -> ApiResult<Json<Vec<PlatformResponse>>> {
    info!("accessed GET /v1/platforms");

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let repo = MySqlMasterRepository(db);

    let platforms = use_case::get_platforms(&repo).await?;
    Ok(Json(platforms))
}

// ─── ルーター ──────────────────────────────────────────────────

/// 認証必要ルート: GET /genres, GET /platforms
pub fn routes() -> Router<AppState> {
    Router::new()
        .route("/genres", get(get_genres))
        .route("/platforms", get(get_platforms))
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use std::sync::Arc;

    use axum::Router;
    use axum::http::StatusCode;
    use axum::middleware;
    use axum_test::TestServer;

    use crate::app::router::AppState;
    use crate::config::{
        Config, DatabaseConfig, JwtConfig, LogConfig, RedisConfig, ServerConfig, TmdbConfig,
    };
    use crate::pkg::jwt;
    use crate::pkg::middleware::auth::require_auth;

    use super::*;

    // ── テスト用ヘルパー ─────────────────────────────────────────

    fn test_config() -> Arc<Config> {
        Arc::new(Config {
            server: ServerConfig { port: 8080 },
            database: DatabaseConfig {
                writer_host: String::new(),
                writer_name: String::new(),
                writer_user: String::new(),
                writer_password: String::new(),
                writer_core_count: 1,
                reader_host: String::new(),
                reader_name: String::new(),
                reader_user: String::new(),
                reader_password: String::new(),
                reader_core_count: 1,
            },
            redis: RedisConfig {
                host: String::new(),
                port: 6379,
                password: String::new(),
                db: 0,
            },
            log: LogConfig {
                level: "info".to_string(),
            },
            jwt: JwtConfig {
                secret: "test_secret".to_string(),
            },
            tmdb: TmdbConfig {
                access_token: String::new(),
            },
        })
    }

    fn state_without_db() -> AppState {
        AppState {
            config: test_config(),
            db: None,
            tmdb: Arc::new(crate::pkg::tmdb::TmdbClient::new("")),
            redis: None,
        }
    }

    fn valid_bearer(user_id: i32) -> String {
        let token = jwt::generate_access_token(user_id, 3600, "test_secret").unwrap();
        format!("Bearer {token}")
    }

    fn test_router(state: AppState) -> axum::Router {
        let config = Arc::clone(&state.config);
        Router::new()
            .merge(routes())
            .layer(middleware::from_fn_with_state(config, require_auth))
            .with_state(state)
    }

    // ── GET /genres ───────────────────────────────────────────────

    /// 認証ヘッダーなしは 401
    #[tokio::test]
    async fn test_get_genres_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let response = server.get("/genres").await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 無効なトークンは 401
    #[tokio::test]
    async fn test_get_genres_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let response = server
            .get("/genres")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 期限切れトークンは 401
    #[tokio::test]
    async fn test_get_genres_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();

        let response = server
            .get("/genres")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// Bearer プレフィックスなしは 401
    #[tokio::test]
    async fn test_get_genres_unauthorized_without_bearer_prefix() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let token = jwt::generate_access_token(1, 3600, "test_secret").unwrap();

        let response = server
            .get("/genres")
            .add_header("Authorization", token)
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 有効なトークンで DB なしは 500（認証は通る）
    #[tokio::test]
    async fn test_get_genres_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let response = server
            .get("/genres")
            .add_header("Authorization", valid_bearer(1))
            .await;

        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── GET /platforms ────────────────────────────────────────────

    /// 認証ヘッダーなしは 401
    #[tokio::test]
    async fn test_get_platforms_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let response = server.get("/platforms").await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 無効なトークンは 401
    #[tokio::test]
    async fn test_get_platforms_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let response = server
            .get("/platforms")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 期限切れトークンは 401
    #[tokio::test]
    async fn test_get_platforms_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();

        let response = server
            .get("/platforms")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// Bearer プレフィックスなしは 401
    #[tokio::test]
    async fn test_get_platforms_unauthorized_without_bearer_prefix() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let token = jwt::generate_access_token(1, 3600, "test_secret").unwrap();

        let response = server
            .get("/platforms")
            .add_header("Authorization", token)
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 有効なトークンで DB なしは 500（認証は通る）
    #[tokio::test]
    async fn test_get_platforms_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();

        let response = server
            .get("/platforms")
            .add_header("Authorization", valid_bearer(1))
            .await;

        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }
}
