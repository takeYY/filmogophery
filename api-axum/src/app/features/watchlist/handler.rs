use axum::{
    Json, Router,
    extract::{Path, Query, State},
    http::StatusCode,
    routing::{delete, get, post},
};
use serde::Deserialize;
use std::collections::HashMap;
use tracing::info;

use crate::app::responses::{ApiResult, AppError};
use crate::app::router::AppState;
use crate::pkg::middleware::auth::AuthUser;

use super::repository::{MySqlMovieExistsRepository, MySqlWatchlistRepository};
use super::use_case::{self, AddToWatchlistInput, WatchlistItemResponse};

// ─── クエリパラメータ ──────────────────────────────────────────

#[derive(Debug, Deserialize)]
pub struct GetWatchlistQuery {
    #[serde(default = "default_limit")]
    pub limit: i32,
    #[serde(default = "default_offset")]
    pub offset: i32,
}

fn default_limit() -> i32 { 12 }
fn default_offset() -> i32 { 0 }

// ─── ハンドラ ──────────────────────────────────────────────────

/// GET /v1/watchlist — ウォッチリスト一覧（認証必要）
async fn get_watchlist(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Query(query): Query<GetWatchlistQuery>,
) -> ApiResult<Json<Vec<WatchlistItemResponse>>> {
    info!("accessed GET /v1/watchlist");

    if query.limit < 1 || query.limit > 12 {
        let mut errors = HashMap::new();
        errors.insert("limit".to_string(), vec!["limit must be between 1 and 12".to_string()]);
        return Err(AppError::ValidationError(errors));
    }
    if query.offset < 0 {
        let mut errors = HashMap::new();
        errors.insert("offset".to_string(), vec!["offset must be 0 or greater".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let repo = MySqlWatchlistRepository(db);

    let items = use_case::get_watchlist(&repo, claims.user_id, query.limit, query.offset).await?;
    Ok(Json(items))
}

/// POST /v1/watchlist — ウォッチリスト登録（認証必要）
async fn add_to_watchlist(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Json(body): Json<AddToWatchlistInput>,
) -> ApiResult<StatusCode> {
    info!("accessed POST /v1/watchlist");

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let watchlist_repo = MySqlWatchlistRepository(db);
    let movie_repo = MySqlMovieExistsRepository(db);

    use_case::add_to_watchlist(&watchlist_repo, &movie_repo, claims.user_id, body).await?;
    Ok(StatusCode::CREATED)
}

/// DELETE /v1/watchlist/{watchlistId} — ウォッチリストから削除（認証必要）
async fn remove_from_watchlist(
    State(state): State<AppState>,
    AuthUser(_claims): AuthUser,
    Path(watchlist_id): Path<i32>,
) -> ApiResult<StatusCode> {
    info!("accessed DELETE /v1/watchlist/{watchlist_id}");

    if watchlist_id < 1 {
        let mut errors = HashMap::new();
        errors.insert("watchlistId".to_string(), vec!["watchlistId must be a positive integer".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let repo = MySqlWatchlistRepository(db);

    use_case::remove_from_watchlist(&repo, watchlist_id).await?;
    Ok(StatusCode::NO_CONTENT)
}

// ─── ルーター ──────────────────────────────────────────────────

pub fn routes() -> Router<AppState> {
    Router::new()
        .route("/watchlist", get(get_watchlist))
        .route("/watchlist", post(add_to_watchlist))
        .route("/watchlist/{watchlistId}", delete(remove_from_watchlist))
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
    use crate::pkg::tmdb::TmdbClient;

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
            redis: RedisConfig { host: String::new(), port: 6379, password: String::new(), db: 0 },
            log: LogConfig { level: "info".to_string() },
            jwt: JwtConfig { secret: "test_secret".to_string() },
            tmdb: TmdbConfig { access_token: String::new() },
        })
    }

    fn state_without_db() -> AppState {
        AppState {
            config: test_config(),
            db: None,
            tmdb: Arc::new(TmdbClient::new("")),
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

    // ── GET /watchlist ────────────────────────────────────────────

    #[tokio::test]
    async fn test_get_watchlist_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server.get("/watchlist").await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_watchlist_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/watchlist")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_watchlist_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .get("/watchlist")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_watchlist_invalid_limit_zero() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/watchlist")
            .add_query_params(&[("limit", "0")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: serde_json::Value = response.json();
        assert!(body["errors"]["limit"].is_array());
    }

    #[tokio::test]
    async fn test_get_watchlist_invalid_limit_over_max() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/watchlist")
            .add_query_params(&[("limit", "13")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_get_watchlist_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/watchlist")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── POST /watchlist ───────────────────────────────────────────

    #[tokio::test]
    async fn test_add_to_watchlist_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/watchlist")
            .json(&serde_json::json!({ "movieId": 1 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_add_to_watchlist_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/watchlist")
            .add_header("Authorization", "Bearer invalid.token.here")
            .json(&serde_json::json!({ "movieId": 1 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_add_to_watchlist_missing_body_returns_415() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/watchlist")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    #[tokio::test]
    async fn test_add_to_watchlist_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/watchlist")
            .add_header("Authorization", valid_bearer(1))
            .json(&serde_json::json!({ "movieId": 1 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── DELETE /watchlist/{watchlistId} ───────────────────────────

    #[tokio::test]
    async fn test_remove_from_watchlist_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server.delete("/watchlist/1").await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_remove_from_watchlist_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .delete("/watchlist/1")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_remove_from_watchlist_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .delete("/watchlist/1")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_remove_from_watchlist_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .delete("/watchlist/1")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }
}
