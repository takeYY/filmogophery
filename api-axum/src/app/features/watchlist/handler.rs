use axum::Router;

use crate::app::router::AppState;

// TODO: GET    /v1/watchlist     (requires auth)
// TODO: POST   /v1/watchlist     (requires auth)
// TODO: DELETE /v1/watchlist/:id (requires auth)

pub fn routes() -> Router<AppState> {
    Router::new()
}
