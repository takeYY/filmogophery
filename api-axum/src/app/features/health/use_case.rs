use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct HealthResponse {
    pub status: String,
}

/// ヘルスチェック use case
pub fn check_health() -> HealthResponse {
    HealthResponse {
        status: "ok".to_string(),
    }
}
