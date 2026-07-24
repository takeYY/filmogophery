use axum::{
    Json, Router,
    extract::{Query, State},
    http::StatusCode,
    routing::{get, post},
};
use serde::Deserialize;
use std::collections::HashMap;
use tracing::info;
use validator::Validate;

use crate::app::responses::{AppError, ApiResult};
use crate::app::router::AppState;
use crate::pkg::middleware::auth::AuthUser;

use super::repository::{MySqlPointRepository, MySqlTokenRepository, MySqlUserRepository};
use super::use_case;

// ─── クエリパラメータ ──────────────────────────────────────────

#[derive(Debug, Deserialize)]
pub struct GetUserPointsQuery {
    #[serde(default = "default_limit")]
    pub limit: i32,
    #[serde(default = "default_offset")]
    pub offset: i32,
}

fn default_limit() -> i32 {
    20
}
fn default_offset() -> i32 {
    0
}

// ─── ハンドラ ──────────────────────────────────────────────────

/// POST /v1/users — ユーザー登録（認証不要）
async fn create_user(
    State(state): State<AppState>,
    Json(body): Json<use_case::CreateUserInput>,
) -> ApiResult<(StatusCode, Json<use_case::TokenResponse>)> {
    info!("accessed POST /v1/users");

    if let Err(errors) = body.validate() {
        let field_errors: HashMap<String, Vec<String>> = errors
            .field_errors()
            .into_iter()
            .map(|(field, errs)| {
                let messages = errs
                    .iter()
                    .filter_map(|e| e.message.as_ref().map(|m| m.to_string()))
                    .collect();
                (field.to_string(), messages)
            })
            .collect();
        return Err(AppError::ValidationError(field_errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let user_repo = MySqlUserRepository(db);
    let token_repo = MySqlTokenRepository(db);

    let token = use_case::create_user(&user_repo, &token_repo, &state.config.jwt.secret, body).await?;
    Ok((StatusCode::CREATED, Json(token)))
}

/// GET /v1/users/me — ログインユーザー取得（認証必要）
async fn get_current_user(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
) -> ApiResult<Json<use_case::UserResponse>> {
    info!("accessed GET /v1/users/me");

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let user_repo = MySqlUserRepository(db);

    let user = use_case::get_current_user(&user_repo, claims.user_id).await?;
    Ok(Json(user))
}

/// GET /v1/users/me/points — ポイント・レベル取得（認証必要）
async fn get_user_points(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
    Query(query): Query<GetUserPointsQuery>,
) -> ApiResult<Json<use_case::UserPointsResponse>> {
    info!("accessed GET /v1/users/me/points");

    if query.limit < 1 || query.limit > 50 {
        let mut errors = HashMap::new();
        errors.insert(
            "limit".to_string(),
            vec!["limit must be between 1 and 50".to_string()],
        );
        return Err(AppError::ValidationError(errors));
    }
    if query.offset < 0 {
        let mut errors = HashMap::new();
        errors.insert(
            "offset".to_string(),
            vec!["offset must be 0 or greater".to_string()],
        );
        return Err(AppError::ValidationError(errors));
    }

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let point_repo = MySqlPointRepository(db);

    let result =
        use_case::get_user_points(&point_repo, claims.user_id, query.limit, query.offset).await?;
    Ok(Json(result))
}

// ─── ルーター ──────────────────────────────────────────────────

pub fn public_routes() -> Router<AppState> {
    Router::new().route("/users", post(create_user))
}

pub fn protected_routes() -> Router<AppState> {
    Router::new()
        .route("/users/me", get(get_current_user))
        .route("/users/me/points", get(get_user_points))
}

// ─── テスト ───────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use std::sync::Arc;

    use axum::http::StatusCode;
    use axum::middleware;
    use axum_test::TestServer;
    use serde_json::Value;

    use crate::app::router::AppState;
    use crate::config::{
        Config, DatabaseConfig, JwtConfig, LogConfig, RedisConfig, ServerConfig, TmdbConfig,
    };
    use crate::pkg::jwt;
    use crate::pkg::middleware::auth::require_auth;

    use super::*;

    // ── テスト用ヘルパー ─────────────────────────────────────────

    /// テスト用の Config を生成する
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

    /// DB なし AppState（public routes のテストに使用）
    fn state_without_db() -> AppState {
        AppState {
            config: test_config(),
            db: None,
            tmdb: Arc::new(crate::pkg::tmdb::TmdbClient::new("")),
            redis: None,
        }
    }

    /// 有効な JWT を生成する（protected routes のテストに使用）
    fn valid_bearer(user_id: i32) -> String {
        let token = jwt::generate_access_token(user_id, 3600, "test_secret").unwrap();
        format!("Bearer {token}")
    }

    /// public_routes のみのルーター
    fn public_test_router(state: AppState) -> axum::Router {
        Router::new()
            .merge(public_routes())
            .with_state(state)
    }

    /// protected_routes に auth ミドルウェアを付けたルーター
    fn protected_test_router(state: AppState) -> axum::Router {
        let config = Arc::clone(&state.config);
        Router::new()
            .merge(protected_routes())
            .layer(middleware::from_fn_with_state(config, require_auth))
            .with_state(state)
    }

    // ── POST /users ──────────────────────────────────────────────

    /// DB なしの場合は 500 を返す（DB 接続が前提のエンドポイント）
    #[tokio::test]
    async fn test_create_user_returns_500_without_db() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/users")
            .json(&serde_json::json!({
                "username": "alice",
                "email": "alice@example.com",
                "password": "password123"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    /// バリデーションエラー: username が空
    #[tokio::test]
    async fn test_create_user_validation_empty_username() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/users")
            .json(&serde_json::json!({
                "username": "",
                "email": "alice@example.com",
                "password": "password123"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: Value = response.json();
        assert!(body["errors"]["username"].is_array());
    }

    /// バリデーションエラー: email の形式が不正
    #[tokio::test]
    async fn test_create_user_validation_invalid_email() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/users")
            .json(&serde_json::json!({
                "username": "alice",
                "email": "not-an-email",
                "password": "password123"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: Value = response.json();
        assert!(body["errors"]["email"].is_array());
    }

    /// バリデーションエラー: password が短すぎる
    #[tokio::test]
    async fn test_create_user_validation_short_password() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/users")
            .json(&serde_json::json!({
                "username": "alice",
                "email": "alice@example.com",
                "password": "short"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: Value = response.json();
        assert!(body["errors"]["password"].is_array());
    }

    // ── GET /users/me ────────────────────────────────────────────

    /// 認証ヘッダーなしは 401
    #[tokio::test]
    async fn test_get_current_user_unauthorized_without_token() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server.get("/users/me").await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 無効なトークンは 401
    #[tokio::test]
    async fn test_get_current_user_unauthorized_with_invalid_token() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server
            .get("/users/me")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 有効なトークンで DB なしは 500（認証は通る）
    #[tokio::test]
    async fn test_get_current_user_returns_500_without_db() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server
            .get("/users/me")
            .add_header("Authorization", valid_bearer(1))
            .await;

        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── GET /users/me/points ──────────────────────────────────────

    /// 認証ヘッダーなしは 401
    #[tokio::test]
    async fn test_get_user_points_unauthorized_without_token() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server.get("/users/me/points").await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// limit が範囲外（0）は 400
    #[tokio::test]
    async fn test_get_user_points_invalid_limit_zero() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server
            .get("/users/me/points")
            .add_query_params(&[("limit", "0")])
            .add_header("Authorization", valid_bearer(1))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: Value = response.json();
        assert!(body["errors"]["limit"].is_array());
    }

    /// limit が範囲外（51）は 400
    #[tokio::test]
    async fn test_get_user_points_invalid_limit_over_max() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server
            .get("/users/me/points")
            .add_query_params(&[("limit", "51")])
            .add_header("Authorization", valid_bearer(1))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
    }

    /// 有効なトークンで DB なしは 500（認証・バリデーションは通る）
    #[tokio::test]
    async fn test_get_user_points_returns_500_without_db() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server
            .get("/users/me/points")
            .add_header("Authorization", valid_bearer(1))
            .await;

        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }
}
