use std::sync::Arc;

use axum::{
    extract::{Request, State},
    middleware::Next,
    response::Response,
};

use crate::app::responses::AppError;
use crate::config::Config;
use crate::pkg::jwt;

/// JWT 認証ミドルウェア
/// Authorization: Bearer <token> ヘッダーを検証し、
/// クレームを Request の Extension に付与する。
pub async fn require_auth(
    State(config): State<Arc<Config>>,
    mut request: Request,
    next: Next,
) -> Result<Response, AppError> {
    let auth_header = request
        .headers()
        .get("Authorization")
        .and_then(|v| v.to_str().ok())
        .ok_or(AppError::Unauthorized)?;

    let token = auth_header
        .strip_prefix("Bearer ")
        .ok_or(AppError::Unauthorized)?;

    let claims = jwt::verify_access_token(token, &config.jwt.secret)?;

    // TODO: DB からユーザーを取得して Extension に付与する
    // 現在はクレームのみ付与
    request.extensions_mut().insert(claims);

    Ok(next.run(request).await)
}
