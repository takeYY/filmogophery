use axum::{Json, Router, routing::get};
use tracing::info;

use crate::app::router::AppState;

use super::use_case;

/// GET /health
async fn check_health() -> Json<use_case::HealthResponse> {
    info!("accessed GET health");
    Json(use_case::check_health())
}

pub fn routes() -> Router<AppState> {
    Router::new().route("/health", get(check_health))
}

#[cfg(test)]
mod tests {
    use axum::http::StatusCode;
    use axum_test::TestServer;

    use super::*;

    fn test_router() -> axum::Router {
        // AppState なしで直接構築し、with_state(()) で解決
        axum::Router::new()
            .route("/health", get(check_health))
            .with_state(())
    }

    #[tokio::test]
    async fn test_check_health_returns_200() {
        let server = TestServer::new(test_router()).unwrap();

        let response = server.get("/health").await;

        assert_eq!(response.status_code(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_check_health_returns_ok_status() {
        let server = TestServer::new(test_router()).unwrap();

        let response = server.get("/health").await;
        let body: serde_json::Value = response.json();

        assert_eq!(body["status"], "ok");
    }
}
