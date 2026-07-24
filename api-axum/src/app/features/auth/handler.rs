use std::collections::HashMap;

use axum::{Json, Router, extract::State, http::StatusCode, routing::post};
use tracing::info;
use validator::Validate;

use crate::app::responses::{ApiResult, AppError};
use crate::app::router::AppState;
use crate::pkg::middleware::auth::AuthUser;

use super::use_case::{self, LoginInput};
use crate::app::features::user::repository::{MySqlTokenRepository, MySqlUserRepository};

// ─── ハンドラ ──────────────────────────────────────────────────

/// POST /v1/auth/login — ログイン（認証不要）
async fn login(
    State(state): State<AppState>,
    Json(body): Json<LoginInput>,
) -> ApiResult<Json<use_case::TokenResponse>> {
    info!("accessed POST /v1/auth/login");

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

    let token = use_case::login(&user_repo, &token_repo, &state.config.jwt.secret, body).await?;
    Ok(Json(token))
}

/// POST /v1/auth/logout — ログアウト（認証必要）
async fn logout(
    State(state): State<AppState>,
    AuthUser(claims): AuthUser,
) -> ApiResult<StatusCode> {
    info!("accessed POST /v1/auth/logout");

    let db = state.db.as_ref().ok_or(AppError::InternalServerError)?;
    let token_repo = MySqlTokenRepository(db);

    use_case::logout(&token_repo, claims.user_id).await?;
    Ok(StatusCode::NO_CONTENT)
}

// ─── ルーター ──────────────────────────────────────────────────

/// 認証不要ルート: POST /auth/login
pub fn public_routes() -> Router<AppState> {
    Router::new().route("/auth/login", post(login))
}

/// 認証必要ルート: POST /auth/logout
pub fn protected_routes() -> Router<AppState> {
    Router::new().route("/auth/logout", post(logout))
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

    fn public_test_router(state: AppState) -> axum::Router {
        Router::new()
            .merge(public_routes())
            .with_state(state)
    }

    fn protected_test_router(state: AppState) -> axum::Router {
        let config = Arc::clone(&state.config);
        Router::new()
            .merge(protected_routes())
            .layer(middleware::from_fn_with_state(config, require_auth))
            .with_state(state)
    }

    // ── POST /auth/login ──────────────────────────────────────────

    /// バリデーションエラー: email が不正
    #[tokio::test]
    async fn test_login_validation_invalid_email() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/auth/login")
            .json(&serde_json::json!({
                "email": "not-an-email",
                "password": "password123"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: serde_json::Value = response.json();
        assert!(body["errors"]["email"].is_array());
    }

    /// バリデーションエラー: password が短すぎる
    #[tokio::test]
    async fn test_login_validation_short_password() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/auth/login")
            .json(&serde_json::json!({
                "email": "test@example.com",
                "password": "short"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: serde_json::Value = response.json();
        assert!(body["errors"]["password"].is_array());
    }

    /// バリデーションエラー: email・password 両方が不正 → errors に両フィールドが含まれる
    #[tokio::test]
    async fn test_login_validation_both_fields_invalid() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/auth/login")
            .json(&serde_json::json!({
                "email": "not-an-email",
                "password": "short"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::BAD_REQUEST);
        let body: serde_json::Value = response.json();
        assert!(body["errors"]["email"].is_array());
        assert!(body["errors"]["password"].is_array());
    }

    /// リクエストボディなし → 415 (axum が Content-Type なしで JSON extractor に失敗)
    #[tokio::test]
    async fn test_login_missing_body_returns_415() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server.post("/auth/login").await;

        assert_eq!(response.status_code(), StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    /// DB なしの場合は 500（バリデーション通過後に DB アクセス）
    #[tokio::test]
    async fn test_login_returns_500_without_db() {
        let server = TestServer::new(public_test_router(state_without_db())).unwrap();

        let response = server
            .post("/auth/login")
            .json(&serde_json::json!({
                "email": "test@example.com",
                "password": "password123"
            }))
            .await;

        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }

    // ── POST /auth/logout ─────────────────────────────────────────

    /// 認証ヘッダーなしは 401
    #[tokio::test]
    async fn test_logout_unauthorized_without_token() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server.post("/auth/logout").await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 無効なトークンは 401
    #[tokio::test]
    async fn test_logout_unauthorized_with_invalid_token() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server
            .post("/auth/logout")
            .add_header("Authorization", "Bearer invalid.token.here")
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 期限切れトークンは 401
    #[tokio::test]
    async fn test_logout_unauthorized_with_expired_token() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        // leeway (デフォルト60秒) を超える過去のトークンを生成
        let token = jwt::generate_access_token(1, -3700, "test_secret").unwrap();

        let response = server
            .post("/auth/logout")
            .add_header("Authorization", format!("Bearer {token}"))
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// Bearer プレフィックスなしは 401
    #[tokio::test]
    async fn test_logout_unauthorized_without_bearer_prefix() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let token = jwt::generate_access_token(1, 3600, "test_secret").unwrap();

        let response = server
            .post("/auth/logout")
            .add_header("Authorization", token) // "Bearer " なし
            .await;

        assert_eq!(response.status_code(), StatusCode::UNAUTHORIZED);
    }

    /// 有効なトークンで DB なしは 500（認証は通る）
    #[tokio::test]
    async fn test_logout_returns_500_without_db() {
        let server = TestServer::new(protected_test_router(state_without_db())).unwrap();

        let response = server
            .post("/auth/logout")
            .add_header("Authorization", valid_bearer(1))
            .await;

        assert_eq!(response.status_code(), StatusCode::INTERNAL_SERVER_ERROR);
    }
}
