use axum::{
    Router,
    extract::{Path, State},
    http::StatusCode,
    routing::{post, put},
    Json,
};
use tracing::info;

use crate::app::responses::{ApiResult, AppError};
use crate::app::router::AppState;
use crate::pkg::middleware::auth::AuthUser;

use super::repository::{
    MySqlMovieExistsRepository, MySqlPlatformExistsRepository,
    MySqlPointRepository, MySqlReviewRepository, MySqlWatchHistoryRepository,
};
use super::use_case::{self, CreateReviewInput, UpdateReviewInput};

// ─── ハンドラ ──────────────────────────────────────────────────

/// POST /v1/movies/{movieId}/reviews — レビュー登録（認証必要）
async fn create_review(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Path(movie_id): Path<i32>,
    Json(body): Json<CreateReviewInput>,
) -> ApiResult<StatusCode> {
    info!("accessed POST /v1/movies/{movie_id}/reviews");

    if movie_id < 1 {
        let mut errors = std::collections::HashMap::new();
        errors.insert("movieId".to_string(), vec!["movieId must be a positive integer".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;

    let review_repo = MySqlReviewRepository(db);
    let movie_repo = MySqlMovieExistsRepository(db);
    let platform_repo = MySqlPlatformExistsRepository(db);
    let wh_repo = MySqlWatchHistoryRepository(db);
    let point_repo = MySqlPointRepository(db);

    use_case::create_review(
        &review_repo,
        &movie_repo,
        &platform_repo,
        &wh_repo,
        &point_repo,
        claims.user_id,
        movie_id,
        body,
    )
    .await?;

    Ok(StatusCode::CREATED)
}

/// PUT /v1/reviews/{reviewId} — レビュー更新（認証必要）
async fn update_review(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Path(review_id): Path<i32>,
    Json(body): Json<UpdateReviewInput>,
) -> ApiResult<StatusCode> {
    info!("accessed PUT /v1/reviews/{review_id}");

    if review_id < 1 {
        let mut errors = std::collections::HashMap::new();
        errors.insert("reviewId".to_string(), vec!["reviewId must be a positive integer".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let review_repo = MySqlReviewRepository(db);

    use_case::update_review(&review_repo, claims.user_id, review_id, body).await?;

    Ok(StatusCode::NO_CONTENT)
}

// ─── ルーター ──────────────────────────────────────────────────

pub fn routes() -> Router<AppState> {
    Router::new()
        .route("/movies/{movieId}/reviews", post(create_review))
        .route("/reviews/{reviewId}", put(update_review))
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
            redis: RedisConfig {
                host: String::new(),
                port: 6379,
                password: String::new(),
                db: 0,
            },
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

    // ── POST /movies/{movieId}/reviews ────────────────────────────

    #[tokio::test]
    async fn test_create_review_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/movies/1/reviews")
            .json(&serde_json::json!({ "rating": 4.0 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_create_review_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/movies/1/reviews")
            .add_header("Authorization", "Bearer invalid.token")
            .json(&serde_json::json!({ "rating": 4.0 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_create_review_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .post("/movies/1/reviews")
            .add_header("Authorization", format!("Bearer {token}"))
            .json(&serde_json::json!({ "rating": 4.0 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_create_review_missing_body_returns_415() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/movies/1/reviews")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    #[tokio::test]
    async fn test_create_review_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/movies/1/reviews")
            .add_header("Authorization", valid_bearer(1))
            .json(&serde_json::json!({ "rating": 4.0 }))
            .await;
        // validation は use_case で行うが DB なし → 500
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── PUT /reviews/{reviewId} ───────────────────────────────────

    #[tokio::test]
    async fn test_update_review_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .put("/reviews/1")
            .json(&serde_json::json!({ "rating": 5.0 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_update_review_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .put("/reviews/1")
            .add_header("Authorization", "Bearer invalid.token")
            .json(&serde_json::json!({ "rating": 5.0 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_update_review_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .put("/reviews/1")
            .add_header("Authorization", format!("Bearer {token}"))
            .json(&serde_json::json!({ "rating": 5.0 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_update_review_missing_body_returns_415() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .put("/reviews/1")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    #[tokio::test]
    async fn test_update_review_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .put("/reviews/1")
            .add_header("Authorization", valid_bearer(1))
            .json(&serde_json::json!({ "rating": 5.0 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }
}
