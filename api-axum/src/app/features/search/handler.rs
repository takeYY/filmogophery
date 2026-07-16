use axum::Router;

use crate::app::router::AppState;

// TODO: GET /v1/search (requires auth)

pub fn routes() -> Router<AppState> {
    Router::new()
}
