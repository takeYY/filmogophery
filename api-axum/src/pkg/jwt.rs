use chrono::Utc;
use jsonwebtoken::{DecodingKey, EncodingKey, Header, Validation, decode, encode};
use serde::{Deserialize, Serialize};

use crate::app::responses::AppError;

/// JWT クレーム
/// Echo/Hono と同じ構造: { "user_id": i32, "exp": i64, "iat": i64 }
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Claims {
    pub user_id: i32,
    pub exp: i64,
    pub iat: i64,
}

/// Access Token を生成する（HS256）
pub fn generate_access_token(
    user_id: i32,
    expires_in_seconds: i64,
    secret: &str,
) -> Result<String, AppError> {
    let now = Utc::now().timestamp();
    let claims = Claims {
        user_id,
        exp: now + expires_in_seconds,
        iat: now,
    };

    encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(secret.as_bytes()),
    )
    .map_err(|_| AppError::InternalServerError)
}

/// Access Token を検証してクレームを返す
pub fn verify_access_token(token: &str, secret: &str) -> Result<Claims, AppError> {
    decode::<Claims>(
        token,
        &DecodingKey::from_secret(secret.as_bytes()),
        &Validation::default(),
    )
    .map(|data| data.claims)
    .map_err(|_| AppError::Unauthorized)
}
