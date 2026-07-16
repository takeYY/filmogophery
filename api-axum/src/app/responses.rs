use axum::{
    Json,
    http::StatusCode,
    response::{IntoResponse, Response},
};
use serde::Serialize;
use std::collections::HashMap;

/// エラーレスポンスの共通型
/// Echo/Hono と同じ形式: { "message": "...", "errors": { "field": ["msg"] } }
#[derive(Debug, Serialize)]
pub struct ErrorResponse {
    pub message: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub errors: Option<HashMap<String, Vec<String>>>,
}

/// アプリケーション共通エラー型
#[derive(Debug, thiserror::Error)]
pub enum AppError {
    #[error("validation failed")]
    ValidationError(HashMap<String, Vec<String>>),

    #[error("bad request")]
    BadRequest(HashMap<String, Vec<String>>),

    #[error("unauthorized")]
    Unauthorized,

    #[error("{0} not found")]
    NotFound(String),

    #[error("{0} is already exist")]
    Conflict(String),

    #[error("system error")]
    InternalServerError,
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, body) = match self {
            AppError::ValidationError(errors) => (
                StatusCode::BAD_REQUEST,
                ErrorResponse {
                    message: "validation failed".to_string(),
                    errors: Some(errors),
                },
            ),
            AppError::BadRequest(errors) => (
                StatusCode::BAD_REQUEST,
                ErrorResponse {
                    message: "bad request".to_string(),
                    errors: Some(errors),
                },
            ),
            AppError::Unauthorized => (
                StatusCode::UNAUTHORIZED,
                ErrorResponse {
                    message: "unauthorized".to_string(),
                    errors: None,
                },
            ),
            AppError::NotFound(resource) => (
                StatusCode::NOT_FOUND,
                ErrorResponse {
                    message: format!("{} not found", resource),
                    errors: None,
                },
            ),
            AppError::Conflict(resource) => (
                StatusCode::CONFLICT,
                ErrorResponse {
                    message: format!("{} is already exist", resource),
                    errors: None,
                },
            ),
            AppError::InternalServerError => (
                StatusCode::INTERNAL_SERVER_ERROR,
                ErrorResponse {
                    message: "system error".to_string(),
                    errors: None,
                },
            ),
        };

        (status, Json(body)).into_response()
    }
}

pub type ApiResult<T> = Result<T, AppError>;
