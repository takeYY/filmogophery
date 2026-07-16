use axum::Router;

use crate::app::router::AppState;

// TODO: GET /v1/master/genres    (public)
// TODO: GET /v1/master/platforms (public)

pub fn routes() -> Router<AppState> {
    Router::new()
}
