use axum::{
    Json, Router,
    extract::{Path, Query, State},
    http::StatusCode,
    routing::{get, post},
};
use serde::Deserialize;
use std::collections::HashMap;
use tracing::info;

use crate::app::responses::{ApiResult, AppError};
use crate::app::router::AppState;
use crate::pkg::middleware::auth::AuthUser;

use super::repository::MySqlMovieRepository;
use super::repository::MovieRepository as _;
use super::use_case::{self, MovieDetailResponse, MovieResponse, MovieWatchHistoryResponse, WatchHistoryResponse};
// ─── クエリパラメータ ──────────────────────────────────────────

#[derive(Debug, Deserialize)]
pub struct GetMoviesQuery {
    pub genre: Option<String>,
    #[serde(default = "default_movie_limit")]
    pub limit: i32,
    #[serde(default = "default_offset")]
    pub offset: i32,
}

#[derive(Debug, Deserialize)]
pub struct SearchMoviesQuery {
    pub title: Option<String>,
    #[serde(default = "default_search_limit")]
    pub limit: i32,
    #[serde(default = "default_offset")]
    pub offset: i32,
}

#[derive(Debug, Deserialize)]
pub struct GetWatchHistoryQuery {
    #[serde(default = "default_movie_limit")]
    pub limit: i32,
    #[serde(default = "default_offset")]
    pub offset: i32,
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CreateWatchHistoryInput {
    pub platform_id: i32,
    pub watched_date: Option<String>, // "YYYY-MM-DD"
}

fn default_movie_limit() -> i32 { 12 }
fn default_search_limit() -> i32 { 20 }
fn default_offset() -> i32 { 0 }

/// 上映時間からポイントを計算する（Echo/Hono と同じロジック）
fn calc_watch_points(runtime_minutes: i32) -> i32 {
    if runtime_minutes <= 90 { 10 }
    else if runtime_minutes <= 150 { 15 }
    else { 20 }
}

// ─── ハンドラ ──────────────────────────────────────────────────

/// GET /v1/movies — レビュー済み映画一覧（認証必要）
async fn get_movies(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Query(query): Query<GetMoviesQuery>,
) -> ApiResult<Json<Vec<MovieResponse>>> {
    info!("accessed GET /v1/movies");

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
    let repo = MySqlMovieRepository(db);

    let movies = use_case::get_movies(
        &repo,
        claims.user_id,
        query.genre.as_deref(),
        query.limit,
        query.offset,
    )
    .await?;

    Ok(Json(movies))
}

/// GET /v1/movies/:id — 映画詳細（認証必要）
async fn get_movie_detail(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Path(id): Path<i32>,
) -> ApiResult<Json<MovieDetailResponse>> {
    let movie_id = id;
    info!("accessed GET /v1/movies/{movie_id}");

    if movie_id < 1 {
        let mut errors = HashMap::new();
        errors.insert("id".to_string(), vec!["id must be a positive integer".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let repo = MySqlMovieRepository(db);

    let detail = use_case::get_movie_detail(&repo, &state.tmdb, movie_id, claims.user_id).await?;
    Ok(Json(detail))
}

/// GET /v1/search/movies — 映画検索（認証必要）
async fn search_movies(
    State(state): State<AppState>,
    AuthUser(_claims): AuthUser,
    Query(query): Query<SearchMoviesQuery>,
) -> ApiResult<Json<Vec<MovieResponse>>> {
    info!("accessed GET /v1/search/movies");

    let title = query.title.as_deref().unwrap_or("").trim().to_string();
    if title.is_empty() {
        let mut errors = HashMap::new();
        errors.insert("title".to_string(), vec!["title is required".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    if query.limit < 1 || query.limit > 20 {
        let mut errors = HashMap::new();
        errors.insert("limit".to_string(), vec!["limit must be between 1 and 20".to_string()]);
        return Err(AppError::ValidationError(errors));
    }
    if query.offset < 0 {
        let mut errors = HashMap::new();
        errors.insert("offset".to_string(), vec!["offset must be 0 or greater".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let repo = MySqlMovieRepository(db);

    let movies = use_case::search_movies(
        &repo,
        &state.tmdb,
        state.redis.as_deref(),
        &title,
        query.limit,
        query.offset,
    )
    .await?;
    Ok(Json(movies))
}

/// GET /v1/movies/{movieId}/watch-history — 映画の視聴履歴一覧（認証必要）
async fn get_movie_watch_history(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Path(movie_id): Path<i32>,
) -> ApiResult<Json<Vec<MovieWatchHistoryResponse>>> {
    info!("accessed GET /v1/movies/{movie_id}/watch-history");

    if movie_id < 1 {
        let mut errors = HashMap::new();
        errors.insert("movieId".to_string(), vec!["movieId must be a positive integer".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let repo = MySqlMovieRepository(db);

    // 映画存在確認
    if !repo.movie_exists(movie_id).await? {
        return Err(AppError::NotFound("movie".to_string()));
    }

    let result = use_case::get_movie_watch_history(&repo, claims.user_id, movie_id).await?;
    Ok(Json(result))
}

/// POST /v1/movies/{movieId}/watch-history — 映画の視聴履歴追加（認証必要）
async fn create_movie_watch_history(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Path(movie_id): Path<i32>,
    Json(body): Json<CreateWatchHistoryInput>,
) -> ApiResult<StatusCode> {
    info!("accessed POST /v1/movies/{movie_id}/watch-history");

    if movie_id < 1 {
        let mut errors = HashMap::new();
        errors.insert("movieId".to_string(), vec!["movieId must be a positive integer".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    if body.platform_id < 1 {
        let mut errors = HashMap::new();
        errors.insert("platformId".to_string(), vec!["platformId must be a positive integer".to_string()]);
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;

    // 映画の存在確認（runtime_minutes も取得）
    use crate::app::features::review::repository::{
        MovieExistsRepository, MySqlMovieExistsRepository, MySqlPlatformExistsRepository,
        MySqlPointRepository, MySqlWatchHistoryRepository, PlatformExistsRepository,
        PointRepository, WatchHistoryRepository,
    };
    let movie_repo = MySqlMovieExistsRepository(db);
    let movie = movie_repo
        .find_by_id(movie_id)
        .await?
        .ok_or_else(|| AppError::NotFound("movie".to_string()))?;

    // プラットフォームの存在確認
    let platform_repo = MySqlPlatformExistsRepository(db);
    if !platform_repo.exists(body.platform_id).await? {
        return Err(AppError::NotFound("platform".to_string()));
    }

    // 視聴履歴を作成
    let wh_repo = MySqlWatchHistoryRepository(db);
    let wh_id = wh_repo
        .create(
            claims.user_id,
            movie_id,
            body.platform_id,
            body.watched_date.as_deref(),
        )
        .await?;

    // ポイント付与（上映時間による段階的ポイント）
    let watch_points = calc_watch_points(movie.runtime_minutes);
    let point_repo = MySqlPointRepository(db);
    point_repo
        .grant_points(claims.user_id, watch_points, "watch_history", wh_id)
        .await?;

    Ok(StatusCode::CREATED)
}

/// GET /v1/users/me/watch-history — ユーザーの視聴履歴一覧（認証必要）
async fn get_watch_history(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Query(query): Query<GetWatchHistoryQuery>,
) -> ApiResult<Json<Vec<WatchHistoryResponse>>> {
    info!("accessed GET /v1/users/me/watch-history");

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
    let repo = MySqlMovieRepository(db);

    let result = use_case::get_watch_history(&repo, claims.user_id, query.limit, query.offset).await?;
    Ok(Json(result))
}

// ─── ルーター ──────────────────────────────────────────────────

pub fn routes() -> Router<AppState> {
    Router::new()
        .route("/movies", get(get_movies))
        .route("/movies/{id}", get(get_movie_detail))
        .route("/movies/{movieId}/watch-history", get(get_movie_watch_history))
        .route("/movies/{movieId}/watch-history", post(create_movie_watch_history))
        .route("/search/movies", get(search_movies))
        .route("/users/me/watch-history", get(get_watch_history))
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

    // ── GET /movies ───────────────────────────────────────────────

    #[tokio::test]
    async fn test_get_movies_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server.get("/movies").await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movies_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movies_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .get("/movies")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movies_invalid_limit_zero() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies")
            .add_query_params(&[("limit", "0")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: serde_json::Value = response.json();
        assert!(body["errors"]["limit"].is_array());
    }

    #[tokio::test]
    async fn test_get_movies_invalid_limit_over_max() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies")
            .add_query_params(&[("limit", "13")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_get_movies_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── GET /movies/:id ───────────────────────────────────────────

    #[tokio::test]
    async fn test_get_movie_detail_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server.get("/movies/1").await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movie_detail_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies/1")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movie_detail_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies/1")
            .add_header("Authorization", valid_bearer(1))
            .await;
        // DB なし → InternalServerError
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── GET /search/movies ────────────────────────────────────────

    #[tokio::test]
    async fn test_search_movies_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server.get("/search/movies").await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_search_movies_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/search/movies")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_search_movies_returns_400_when_title_missing() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/search/movies")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: serde_json::Value = response.json();
        assert!(body["errors"]["title"].is_array());
    }

    #[tokio::test]
    async fn test_search_movies_returns_400_when_title_empty() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/search/movies")
            .add_query_params(&[("title", "")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_search_movies_returns_400_when_limit_invalid() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/search/movies")
            .add_query_params(&[("title", "test"), ("limit", "0")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_search_movies_returns_500_without_db() {
        // DB なし → TMDB 呼び出し後に DB アクセスで 500
        // (TMDB アクセストークンが空なので TMDB 呼び出し自体が 500 になることも許容)
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/search/movies")
            .add_query_params(&[("title", "inception")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── GET /movies/{movieId}/watch-history ───────────────────────

    #[tokio::test]
    async fn test_get_movie_watch_history_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server.get("/movies/1/watch-history").await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movie_watch_history_unauthorized_with_invalid_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies/1/watch-history")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movie_watch_history_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .get("/movies/1/watch-history")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_movie_watch_history_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/movies/1/watch-history")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── POST /movies/{movieId}/watch-history ──────────────────────

    #[tokio::test]
    async fn test_create_movie_watch_history_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/movies/1/watch-history")
            .json(&serde_json::json!({ "platformId": 1 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_create_movie_watch_history_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .post("/movies/1/watch-history")
            .add_header("Authorization", format!("Bearer {token}"))
            .json(&serde_json::json!({ "platformId": 1 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_create_movie_watch_history_missing_body_returns_415() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/movies/1/watch-history")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    #[tokio::test]
    async fn test_create_movie_watch_history_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .post("/movies/1/watch-history")
            .add_header("Authorization", valid_bearer(1))
            .json(&serde_json::json!({ "platformId": 1 }))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── GET /users/me/watch-history ───────────────────────────────

    #[tokio::test]
    async fn test_get_watch_history_unauthorized_without_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server.get("/users/me/watch-history").await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_watch_history_unauthorized_with_expired_token() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();
        let response = server
            .get("/users/me/watch-history")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;
        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    #[tokio::test]
    async fn test_get_watch_history_invalid_limit_zero() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/users/me/watch-history")
            .add_query_params(&[("limit", "0")])
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: serde_json::Value = response.json();
        assert!(body["errors"]["limit"].is_array());
    }

    #[tokio::test]
    async fn test_get_watch_history_returns_500_without_db() {
        let server = TestServer::new(test_router(state_without_db())).unwrap();
        let response = server
            .get("/users/me/watch-history")
            .add_header("Authorization", valid_bearer(1))
            .await;
        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }
}
