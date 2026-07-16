use axum::Router;

use crate::app::router::AppState;

// TODO: GET /v1/trending (requires auth)

pub fn routes() -> Router<AppState> {
    Router::new()
}
