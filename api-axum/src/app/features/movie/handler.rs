use axum::Router;

use crate::app::router::AppState;

// TODO: GET  /v1/movies               (requires auth)
// TODO: GET  /v1/movies/:id           (requires auth)
// TODO: GET  /v1/movies/:id/watch-history (requires auth)

pub fn routes() -> Router<AppState> {
    Router::new()
}
